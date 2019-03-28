package server

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type Config struct {
	Bind string

	MongoUser     string
	MongoPassword string
	MongoDb       string
	MongoHost     string
}

func (c *Config) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("PortServerConfig", pflag.PanicOnError)

	//grpc
	f.StringVar(&c.Bind, "bind", "127.0.0.1:0", "ip:port")
	// mongo
	f.StringVar(&c.MongoUser, "mongo_user", "", "mongo db user")
	f.StringVar(&c.MongoPassword, "mongo_password", "", "mongo db password")
	f.StringVar(&c.MongoDb, "mongo_db", "", "mongo database name")
	f.StringVar(&c.MongoDb, "mongo_host", "", "mongo host")

	return f
}

func (c *Config) Validate() error {
	if c.MongoUser == "" || c.MongoPassword == "" || c.MongoDb == "" || c.MongoHost == "" {
		return errors.New(fmt.Sprintf("invlaid mongo config. cfg:%s:%s-%s/%s", c.MongoUser, c.MongoPassword, c.MongoHost, c.MongoDb))
	}
	return nil
}
