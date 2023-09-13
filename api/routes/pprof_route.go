package routes

import (
	"github.com/labstack/echo/v4"
	"manuel71sj/go-api-template/lib"
	"net/http"
	"net/http/pprof"
)

type PprofRoutes struct {
	logger  lib.Logger
	handler lib.HttpHandler
}

// Setup pprof routes
func (r PprofRoutes) Setup() {
	r.logger.Zap.Info("Setting up pprof routes")

	api := r.handler.Engine.Group("/pprof")
	{
		api.GET("/", handler(pprof.Index))
		api.GET("/allocs", handler(pprof.Handler("allocs").ServeHTTP))
		api.GET("/block", handler(pprof.Handler("block").ServeHTTP))
		api.GET("/cmdline", handler(pprof.Cmdline))
		api.GET("/goroutine", handler(pprof.Handler("goroutine").ServeHTTP))
		api.GET("/heap", handler(pprof.Handler("heap").ServeHTTP))
		api.GET("/mutex", handler(pprof.Handler("mutex").ServeHTTP))
		api.GET("/profile", handler(pprof.Profile))
		api.POST("/symbol", handler(pprof.Symbol))
		api.GET("/symbol", handler(pprof.Symbol))
		api.GET("/threadcreate", handler(pprof.Handler("threadcreate").ServeHTTP))
		api.GET("/trace", handler(pprof.Trace))
	}
}

func handler(h http.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}
}

// NewPprofRoutes creates new pprof routes
func NewPprofRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
) PprofRoutes {
	return PprofRoutes{
		logger:  logger,
		handler: handler,
	}
}
