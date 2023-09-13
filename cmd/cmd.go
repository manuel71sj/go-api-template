package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"manuel71sj/go-api-template/cmd/migrate"
	"manuel71sj/go-api-template/cmd/runserver"
	"manuel71sj/go-api-template/cmd/setup"
	"os"
)

var rootCmd = &cobra.Command{
	Use:          "api-backend",
	Short:        "api-backend is a REST API for a application",
	Long:         "Template for a REST API in Go",
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one argument you can view the available parameters through '--help'")
		}

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error { return nil },
	Run:                func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(runserver.StartCmd)
	rootCmd.AddCommand(migrate.StartCmd)
	rootCmd.AddCommand(setup.StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
