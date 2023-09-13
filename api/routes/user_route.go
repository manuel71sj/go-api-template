package routes

import (
	"manuel71sj/go-api-template/api/controllers"
	"manuel71sj/go-api-template/lib"
)

type UserRoutes struct {
	logger         lib.Logger
	handler        lib.HttpHandler
	userController controllers.UserController
}

// Setup user routes
func (r UserRoutes) Setup() {
	r.logger.Zap.Info("Setting up user routes")
	api := r.handler.RouterV1.Group("/users")
	{
		api.GET("", r.userController.Query)
		api.POST("", r.userController.Create)
		api.GET("/:id", r.userController.Get)
		api.PUT("/:id", r.userController.Update)
		api.DELETE("/:id", r.userController.Delete)
		api.POST("/:id/enable", r.userController.Enable)
		api.POST("/:id/disable", r.userController.Disable)
	}
}

// NewUserRoutes creates new user routes
func NewUserRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	userController controllers.UserController,
) UserRoutes {
	return UserRoutes{
		handler:        handler,
		logger:         logger,
		userController: userController,
	}
}
