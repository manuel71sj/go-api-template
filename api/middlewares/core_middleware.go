package middlewares

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"manuel71sj/go-api-template/constants"
	"manuel71sj/go-api-template/lib"
	"runtime"
)

// CoreMiddleware core middleware is a functional extension to "echo",
// including database transactions and panic recovery
// and more
type CoreMiddleware struct {
	handler lib.HttpHandler
	logger  lib.Logger
	db      lib.Database
}

func (m CoreMiddleware) core() echo.MiddlewareFunc {
	logger := m.logger.DesugarZap.With(zap.String("module", "core-mw"))

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			txHandle := m.db.ORM.Begin()
			logger.Info("Beginning database transaction")

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					// recovery stack
					stack := make([]byte, 4<<10)
					length := runtime.Stack(stack, false)
					msg := fmt.Sprintf("PANIC RECOVER: %v%s\n", err, stack[:length])
					logger.Error(msg)

					// rollback database transaction
					logger.Info("Rolling back transaction due to panic")
					txHandle.Rollback()
					ctx.Error(err)
				}
			}()

			ctx.Set(constants.DBTransaction, txHandle)

			if err := next(ctx); err != nil {
				ctx.Error(err)
			}

			code := ctx.Response().Status
			// rollback transaction on server errors
			if code >= 400 {
				msg := fmt.Sprintf("Rolling back transaction due to status code: %d", code)
				logger.Info(msg)

				txHandle.Rollback()
			} else {
				m.logger.DesugarZap.Info("Committing transactions")
				if err := txHandle.Commit().Error; err != nil {
					logger.Error(fmt.Sprintf("Trx commit error: %v", err))
				}
			}

			return nil
		}
	}
}

func (m CoreMiddleware) Setup() {
	m.logger.Zap.Info("Setting up core middleware")
	m.handler.Engine.Use(m.core())
}

// statusInList function checks if context writer status is in provided list
func statusInList(status int, list []int) bool {
	for _, s := range list {
		if s == status {
			return true
		}
	}

	return false
}

// NewCoreMiddleware creates new database transactions middleware
func NewCoreMiddleware(handler lib.HttpHandler, logger lib.Logger, db lib.Database) CoreMiddleware {
	return CoreMiddleware{
		handler: handler,
		logger:  logger,
		db:      db,
	}
}
