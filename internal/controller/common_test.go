package controller

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/ory/dockertest"
	dc "github.com/ory/dockertest/docker"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/port-domain/internal/storage"
	"github.com/port-domain/internal/utils/mongo"
)

var l = logrus.New()

var mgStore globalResource

const mongoDefaultPort = 27017

type globalResource struct {
	mu           sync.Mutex
	resource     interface{}
	shutdownHook func()
}

func TestMain(m *testing.M) {
	var retcode int
	defer func() {
		if mgStore.shutdownHook != nil {
			mgStore.shutdownHook()
		}
		os.Exit(retcode)
	}()
	retcode = m.Run()
}

func NewDockerMongoStorage() (_ *storage.Storage, _ error) {
	mgStore.mu.Lock()
	defer mgStore.mu.Unlock()

	if mgStore.resource != nil {
		return mgStore.resource.(*storage.Storage), nil
	}

	mg, dsCloser, dsError := setupMongoContainer()
	if dsError != nil {
		return nil, dsError
	}
	mgStore.resource = mg
	mgStore.shutdownHook = func() {
		mgStore.mu.Lock()
		defer mgStore.mu.Unlock()
		fmt.Printf("[DEBUG] SHUTTING DOWN MONGODB\n")
		// close mongo
		if err := mg.Client.Disconnect(context.Background()); err != nil {
			l.WithError(err).Error()
			return
		}

		// stop and remove container
		if err := dsCloser(); err != nil {
			l.WithError(err).Error("failed to free docker resources")
		}
		mgStore.resource = nil
		mgStore.shutdownHook = nil
	}

	return mgStore.resource.(*storage.Storage), nil
}

func setupMongoContainer() (_ *storage.Storage, closer func() error, _ error) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	fp, err := freePort()
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	fp = 27017 // todo think how to parametrize
	ePort := []string{fmt.Sprintf("%d", mongoDefaultPort)}
	bPorts := map[dc.Port][]dc.PortBinding{
		dc.Port(fmt.Sprintf("%d", fp)): {{HostPort: fmt.Sprintf("%d", mongoDefaultPort)}},
	}
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "mongo",
		Tag:          "latest",
		ExposedPorts: ePort,
		PortBindings: bPorts,
		Env: []string{
			"MONGO_INITDB_ROOT_USERNAME=test",
			"MONGO_INITDB_ROOT_PASSWORD=root",
			"MONGO_INITDB_DATABASE=admin",
		},
	})
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	if err := pool.Retry(func() error {
		res, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d", fp))
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return errors.New("mongodb is not started yet")
		}

		return nil
	}); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	client, err := mongo.InitAndEnsureMongoDb("test", "root", fmt.Sprintf("localhost:%d", fp), "admin")
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	mongo := storage.New(client)

	var resCloser = func() error {
		return pool.Purge(resource)
	}

	return mongo, resCloser, nil
}

func freePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}
