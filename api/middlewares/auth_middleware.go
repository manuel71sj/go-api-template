package middlewares

import (
	"github.com/labstack/echo/v4"
	"manuel71sj/go-api-template/api/services"
	"manuel71sj/go-api-template/constants"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/pkg/echox"
	"net/http"
	"strings"
)

// AuthMiddleware middleware for cors
type AuthMiddleware struct {
	config      lib.Config
	handler     lib.HttpHandler
	logger      lib.Logger
	authService services.AuthService
}

func (m AuthMiddleware) core() echo.MiddlewareFunc {
	prefixes := m.config.Auth.IgnorePathPrefixes

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			request := ctx.Request()

			if isIgnorePath(request.URL.Path, prefixes...) {
				return next(ctx)
			}

			var (
				auth   = request.Header.Get("Authorization")
				prefix = "Bearer "
				token  string
			)

			if auth != "" && strings.HasPrefix(auth, prefix) {
				token = auth[len(prefix):]
			}

			claims, err := m.authService.ParseToken(token)
			if err != nil {
				return echox.Response{Code: http.StatusUnauthorized, Message: err}.JSON(ctx)
			}

			ctx.Set(constants.CurrentUser, claims)
			return next(ctx)
		}
	}
}

func (m AuthMiddleware) Setup() {
	if !m.config.Auth.Enable {
		return
	}

	m.logger.Zap.Info("Setting up auth middleware")
	m.handler.Engine.Use(m.core())
}

// NewAuthMiddleware creates new cors middleware
func NewAuthMiddleware(
	config lib.Config,
	handler lib.HttpHandler,
	logger lib.Logger,
	authService services.AuthService,
) AuthMiddleware {
	return AuthMiddleware{
		config:      config,
		handler:     handler,
		logger:      logger,
		authService: authService,
	}
}
