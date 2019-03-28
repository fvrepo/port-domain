package utils

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func BindEnv(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		envVar := strings.ToUpper(f.Name)

		if val := os.Getenv(envVar); val != "" {
			cmd.Flags().Set(f.Name, val)
		}
	})
}
