package middlewares

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"manuel71sj/go-api-template/lib"
	"time"
)

// ZapMiddleware middleware for logger
type ZapMiddleware struct {
	handler lib.HttpHandler
	logger  lib.Logger
}

func (m ZapMiddleware) core() echo.MiddlewareFunc {
	logger := m.logger.DesugarZap.With(zap.String("module", "log-mw"))

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()

			if err := next(ctx); err != nil {
				logger = logger.With(zap.Error(err))
				ctx.Error(err)
			}

			request := ctx.Request()
			response := ctx.Response()

			fields := []zapcore.Field{
				zap.String("remote_ip", ctx.RealIP()),
				zap.String("time", time.Since(start).String()),
				zap.String("host", request.Host),
				zap.String("request", fmt.Sprintf("%s %s", request.Method, request.RequestURI)),
				zap.Int("status", response.Status),
				zap.Int64("size", response.Size),
				zap.String("user_agent", request.UserAgent()),
			}

			id := request.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = response.Header().Get(echo.HeaderXRequestID)
				fields = append(fields, zap.String("request_id", id))
			}

			n := response.Status
			switch {
			case n >= 500:
				logger.Error("Server error", fields...)
			case n >= 400:
				logger.Warn("Client error", fields...)
			case n >= 300:
				logger.Info("Redirection", fields...)
			default:
				logger.Info("Success", fields...)
			}

			return nil
		}
	}
}

func (m ZapMiddleware) Setup() {
	m.logger.Zap.Info("Setting up zap middleware")
	m.handler.Engine.Use(m.core())
}

// NewZapMiddleware creates new zap middleware
func NewZapMiddleware(handler lib.HttpHandler, logger lib.Logger) ZapMiddleware {
	return ZapMiddleware{handler: handler, logger: logger}
}
