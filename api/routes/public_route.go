package routes

import (
	"manuel71sj/go-api-template/api/controllers"
	"manuel71sj/go-api-template/lib"
)

type PublicRoutes struct {
	logger            lib.Logger
	handler           lib.HttpHandler
	publicController  controllers.PublicController
	captchaController controllers.CaptchaController
}

// Setup public routes
func (r PublicRoutes) Setup() {
	r.logger.Zap.Info("Setting up public routes")

	api := r.handler.RouterV1.Group("/publics")
	{
		api.GET("/user", r.publicController.UserInfo)
		api.POST("/user/login", r.publicController.UserLogin)
		api.POST("/user/logout", r.publicController.UserLogout)
		api.GET("/user/menutree", r.publicController.MenuTree)

		// sys routes
		api.GET("/sys/routes", r.publicController.SysRoutes)

		// captcha
		api.GET("/captcha", r.captchaController.GetCaptcha)
		api.POST("/captcha/verify", r.captchaController.VerifyCaptcha)
	}
}

// NewPublicRoutes creates new public routes
func NewPublicRoutes(
	logger lib.Logger,
	handler lib.HttpHandler,
	publicController controllers.PublicController,
	captchaController controllers.CaptchaController,
) PublicRoutes {
	return PublicRoutes{
		handler:           handler,
		logger:            logger,
		publicController:  publicController,
		captchaController: captchaController,
	}
}
