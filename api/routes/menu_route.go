package routes

import (
	"manuel71sj/go-api-template/api/controllers"
	"manuel71sj/go-api-template/lib"
)

type MenuRoutes struct {
	logger         lib.Logger
	handler        lib.HttpHandler
	menuController controllers.MenuController
}

// Setup menu routes
func (r MenuRoutes) Setup() {
	r.logger.Zap.Info("Setting up menu routes")

	api := r.handler.RouterV1.Group("/menus")
	{
		api.GET("", r.menuController.Query)

		api.POST("", r.menuController.Create)
		api.GET("/:id", r.menuController.Get)
		api.PUT("/:id", r.menuController.Update)
		api.DELETE("/:id", r.menuController.Delete)
		api.PATCH("/:id/enable", r.menuController.Enable)
		api.PATCH("/:id/disable", r.menuController.Disable)

		api.GET("/:id/actions", r.menuController.GetActions)
		api.PUT("/:id/actions", r.menuController.UpdateActions)
	}
}

// NewMenuRoutes creates new menu routes
func NewMenuRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	menuController controllers.MenuController,
) MenuRoutes {
	return MenuRoutes{
		handler:        handler,
		logger:         logger,
		menuController: menuController,
	}
}
