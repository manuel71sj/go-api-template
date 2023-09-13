package middlewares

import (
	"github.com/labstack/echo/v4/middleware"
	"manuel71sj/go-api-template/lib"
)

// CorsMiddleware middleware for cors
type CorsMiddleware struct {
	handler lib.HttpHandler
	logger  lib.Logger
}

func (m CorsMiddleware) Setup() {
	m.logger.Zap.Info("Setting up cors middleware")

	m.handler.Engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"*"},
	}))
}

// NewCorsMiddleware creates new cors middleware
func NewCorsMiddleware(handler lib.HttpHandler, logger lib.Logger) CorsMiddleware {
	return CorsMiddleware{
		handler: handler,
		logger:  logger,
	}
}
