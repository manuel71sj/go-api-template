package routes

import (
	"manuel71sj/go-api-template/api/controllers"
	"manuel71sj/go-api-template/lib"
)

type RoleRoutes struct {
	logger         lib.Logger
	handler        lib.HttpHandler
	roleController controllers.RoleController
}

// Setup role routes
func (r RoleRoutes) Setup() {
	r.logger.Zap.Info("Setting up role routes")

	api := r.handler.RouterV1.Group("/roles")
	{
		api.GET("", r.roleController.Query)
		api.GET(".all", r.roleController.GetAll)

		api.POST("", r.roleController.Create)
		api.GET("/:id", r.roleController.Get)
		api.PUT("/:id", r.roleController.Update)
		api.DELETE("/:id", r.roleController.Delete)
		api.PATCH("/:id/enable", r.roleController.Enable)
		api.PATCH("/:id/disable", r.roleController.Disable)
	}
}

// NewRoleRoutes creates new role routes
func NewRoleRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	roleController controllers.RoleController,
) RoleRoutes {
	return RoleRoutes{
		handler:        handler,
		logger:         logger,
		roleController: roleController,
	}
}
