package setup

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"manuel71sj/go-api-template/api/repository"
	"manuel71sj/go-api-template/api/services"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models"
	"manuel71sj/go-api-template/pkg/file"
	"os"
)

var (
	configFile string
	menuFile   string

	StartCmd = &cobra.Command{
		Use:          "setup",
		Short:        "Set up data for the application",
		Example:      "{execfile} init -c config/settings.yml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			lib.SetConfigPath(configFile)
		},
		Run: func(cmd *cobra.Command, args []string) {
			config := lib.NewConfig()
			logger := lib.NewLogger(config)
			db := lib.NewDatabase(config, logger)

			menuService := services.NewMenuService(
				logger,
				repository.NewMenuRepository(db, logger),
				repository.NewMenuActionRepository(db, logger),
				repository.NewMenuActionResourceRepository(db, logger),
			)

			if !file.IsFile(menuFile) {
				logger.Zap.Fatal("Menu file does not exist")
			}

			fs, err := os.Open(menuFile)
			if err != nil {
				logger.Zap.Fatalf("Menu file could not be opened: %v", err)
			}

			defer func(fs *os.File) {
				_ = fs.Close()
			}(fs)

			var menuTrees models.MenuTrees
			yd := yaml.NewDecoder(fs)
			if err := yd.Decode(&menuTrees); err != nil {
				logger.Zap.Fatalf("Menu file decode error: %v", err)
			}

			if err := menuService.CreateMenus("", menuTrees); err != nil {
				logger.Zap.Fatalf("Menu file init error: %v", err)
			}

			logger.Zap.Info("Menu file import successfully.")
		},
	}
)

func init() {
	pf := StartCmd.PersistentFlags()
	pf.StringVarP(&configFile, "config", "c",
		"config/config.yaml", "this parameter is used to start the service application.")
	pf.StringVarP(&menuFile, "menu", "m",
		"config/menu.yaml", "this parameter is used to set the initialized menu data.")

	_ = cobra.MarkFlagRequired(pf, "config")
	_ = cobra.MarkFlagRequired(pf, "menu")

}
