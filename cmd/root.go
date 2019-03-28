package cmd

import (
	"github.com/port-domain/cmd/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var l = logrus.New()

var RootCmd = &cobra.Command{
	Use:   "port",
	Short: "Port servce",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		l.WithError(err).Fatal("something goes wrong")
	}
}

func init() {
	RootCmd.AddCommand(server.CMD)
}
