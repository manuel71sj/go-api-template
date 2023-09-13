package routes

import (
	echoSwagger "github.com/swaggo/echo-swagger"
	"manuel71sj/go-api-template/constants"
	"manuel71sj/go-api-template/docs"
	"manuel71sj/go-api-template/lib"
)

// SwaggerRoutes
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http https
// @basePath /api
// @contact.name Manuel Han
// @contact.email manuel71.sj@gmail
type SwaggerRoutes struct {
	config  lib.Config
	logger  lib.Logger
	handler lib.HttpHandler
}

// Setup swagger routes
func (r SwaggerRoutes) Setup() {
	docs.SwaggerInfo.Title = r.config.Name
	docs.SwaggerInfo.Version = constants.Version

	r.logger.Zap.Info("Setting up swagger routes")
	r.handler.Engine.GET("/swagger/*", echoSwagger.WrapHandler)
}

// NewSwaggerRoutes creates new swagger routes
func NewSwaggerRoutes(
	config lib.Config,
	logger lib.Logger,
	handler lib.HttpHandler,
) SwaggerRoutes {
	return SwaggerRoutes{
		config:  config,
		logger:  logger,
		handler: handler,
	}
}
