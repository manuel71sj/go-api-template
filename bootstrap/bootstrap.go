package bootstrap

import (
	"context"
	"go.uber.org/fx"
	"manuel71sj/go-api-template/api/controllers"
	"manuel71sj/go-api-template/api/middlewares"
	"manuel71sj/go-api-template/api/repository"
	"manuel71sj/go-api-template/api/routes"
	"manuel71sj/go-api-template/api/services"
	"manuel71sj/go-api-template/errors"
	"manuel71sj/go-api-template/lib"
	"net/http"
	"time"
)

// Module exported for initializing application
var Module = fx.Options(
	controllers.Module,
	routes.Module,
	lib.Module,
	services.Module,
	middlewares.Module,
	repository.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	handler lib.HttpHandler,
	routes routes.Routes,
	logger lib.Logger,
	config lib.Config,
	middlewares middlewares.Middlewares,
	database lib.Database,
) {
	db, err := database.ORM.DB()
	if err != nil {
		logger.Zap.Fatalf("Error to get database connection: %v", err)
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Zap.Info("Starting application...")

			if err := db.Ping(); err != nil {
				logger.Zap.Fatalf("Error to ping database connection: %v", err)
			}

			// set conn
			db.SetMaxOpenConns(config.Database.MaxOpenConns)
			db.SetMaxIdleConns(config.Database.MaxIdleConns)
			db.SetConnMaxLifetime(time.Duration(config.Database.MaxLifetime) * time.Second)

			go func() {
				middlewares.Setup()
				routes.Setup()

				if err := handler.Engine.Start(config.Http.ListenAddr()); err != nil {
					if errors.Is(err, http.ErrServerClosed) {
						logger.Zap.Debug("Shutting down the Application")
					} else {
						logger.Zap.Fatalf("Error to Start Application: %v", err)
					}
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Zap.Info("Stopping application...")

			_ = handler.Engine.Close()
			_ = db.Close()

			return nil
		},
	})
}
