package middlewares

import (
	"github.com/labstack/echo/v4"
	"manuel71sj/go-api-template/api/services"
	"manuel71sj/go-api-template/constants"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models/dto"
	"manuel71sj/go-api-template/pkg/echox"
	"net/http"
)

type CasbinMiddleware struct {
	handler lib.HttpHandler
	logger  lib.Logger
	config  lib.Config

	casbinService services.CasbinService
}

func (m CasbinMiddleware) core() echo.MiddlewareFunc {
	prefixes := m.config.Casbin.IgnorePathPrefixes

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			request := ctx.Request()

			if isIgnorePath(request.URL.Path, prefixes...) {
				return next(ctx)
			}

			p := ctx.Request().URL.Path
			method := ctx.Request().Method
			claims, ok := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
			if !ok {
				return echox.Response{Code: http.StatusUnauthorized}.JSON(ctx)
			}

			if ok, err := m.casbinService.Enforcer.Enforce(claims.ID, p, method); err != nil {
				return echox.Response{Code: http.StatusForbidden, Message: err}.JSON(ctx)
			} else if !ok {
				return echox.Response{Code: http.StatusForbidden}.JSON(ctx)
			}

			return next(ctx)
		}
	}
}

func (m CasbinMiddleware) Setup() {
	if !m.config.Casbin.Enable {
		return
	}

	m.logger.Zap.Info("Setting up casbin middleware")
	m.handler.Engine.Use(m.core())
}

func NewCasbinMiddleware(
	handler lib.HttpHandler,
	logger lib.Logger,
	config lib.Config,
	casbinService services.CasbinService,
) CasbinMiddleware {
	return CasbinMiddleware{
		handler:       handler,
		logger:        logger,
		config:        config,
		casbinService: casbinService,
	}
}
