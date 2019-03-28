package server

import (
	"math"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/facebookgo/stack"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/port-domain/internal/controller"
	"github.com/port-domain/internal/server"
	"github.com/port-domain/internal/storage"
	"github.com/port-domain/internal/utils"
	"github.com/port-domain/internal/utils/mongo"
	portApi "github.com/port-domain/pkg/grpcapi/port"
)

var cfg Config
var l = logrus.New()

func init() {
	CMD.Flags().AddFlagSet(cfg.Flags())
}

var CMD = &cobra.Command{
	Use:   "server",
	Short: "Ports gRPC service",
	RunE: func(cmd *cobra.Command, args []string) error {
		l.Info("Starting Port gRPC service")
		// parse config from env variables
		utils.BindEnv(cmd)

		if err := cfg.Validate(); err != nil {
			return errors.WithStack(err)
		}

		//init mongodb
		// todo add index for _id
		mongoDb, err := mongo.InitAndEnsureMongoDb(cfg.MongoUser, cfg.MongoPassword, cfg.MongoHost, cfg.MongoDb)
		if err != nil {
			return errors.WithStack(err)
		}
		storage := storage.New(mongoDb)
		controller := controller.New(storage)

		srv := grpc.NewServer(
			grpc.UnaryInterceptor(grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(panicHandler()))),
			grpc.MaxConcurrentStreams(math.MaxUint32),
		)
		portApi.RegisterPortServiceServer(srv, server.New(controller))

		if err := serveWithGracefulStop(srv, cfg); err != nil {
			l.WithError(err).Error("failed to start Ports grpc server")
			return errors.WithStack(err)
		}
		l.Info("Port gRPC service started")

		return nil
	},
}

func panicHandler() grpc_recovery.RecoveryHandlerFunc {
	return func(p interface{}) error {
		l.WithError(errors.New("panic recovery")).
			WithField("stack", stack.Caller(3).String()).
			Error("grpc server panic handler triggered")

		return status.Error(codes.Internal, "Internal error")
	}
}

func serveWithGracefulStop(srv *grpc.Server, config Config) error {
	ln, err := net.Listen("tcp", config.Bind)
	if err != nil {
		l.WithError(err).Fatal("tcp binding failed in grpc server")
	}
	if err := srv.Serve(ln); err != nil {
		return errors.WithStack(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.Signal(syscall.SIGTERM))
	defer func() {
		signal.Stop(c)
		close(c)
	}()

	go func() {
		l.WithField("signal", <-c).Info("got signal")
		srv.GracefulStop()
	}()

	return nil
}
