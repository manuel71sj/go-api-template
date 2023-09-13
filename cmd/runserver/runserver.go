package runserver

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"manuel71sj/go-api-template/bootstrap"
	"manuel71sj/go-api-template/lib"
)

var (
	configFile  string
	casbinModel string

	StartCmd = &cobra.Command{
		Use:          "runserver",
		Short:        "Start API server",
		Example:      "{execfile} server -c config/settings.yml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			lib.SetConfigPath(configFile)
			lib.SetConfigCasbinModelPath(casbinModel)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runApplication()
		},
	}
)

func init() {
	pf := StartCmd.PersistentFlags()
	pf.StringVarP(&configFile, "config", "c",
		"config/config.yaml", "this parameter is used to start the service application.")
	pf.StringVarP(&casbinModel, "casbin", "m",
		"config/casbin_model.conf", "this parameter is used for the running configuration of casbin.")

	_ = cobra.MarkFlagRequired(pf, "config")
	_ = cobra.MarkFlagRequired(pf, "casbin")
}

func runApplication() {
	fx.New(bootstrap.Module, fx.NopLogger).Run()
}
