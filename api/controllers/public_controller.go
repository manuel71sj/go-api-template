package controllers

import (
	"github.com/labstack/echo/v4"
	"manuel71sj/go-api-template/api/services"
	"manuel71sj/go-api-template/constants"
	"manuel71sj/go-api-template/errors"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models/dto"
	"manuel71sj/go-api-template/pkg/echox"
	"net/http"
	"strings"
)

type PublicController struct {
	userService services.UserService
	authService services.AuthService
	captcha     lib.Captcha
	logger      lib.Logger
	config      lib.Config
}

type route struct {
	*echo.Route
	Name *struct{} `json:"name,omitempty"`
}

// SysRoutes
// @Tags Public
// @Summary SysRoutes
// @Produce application/json
// @Success 200 {string} echox.Response "ok"
// @failure 400 {string} echox.Response "bad request"
// @failure 500 {string} echox.Response "internal error"
// @Router /api/v1/publics/sys/routes [get]
func (c PublicController) SysRoutes(ctx echo.Context) error {
	routes := make([]*route, 0)
	for _, eRoute := range ctx.Echo().Routes() {
		// Only interfaces starting with /api/ are exposed
		if !strings.HasPrefix(eRoute.Path, "/api/") {
			continue
		}

		routes = append(routes, &route{Route: eRoute})
	}

	return echox.Response{Code: http.StatusOK, Data: routes}.JSON(ctx)
}

// UserInfo
// @Tags Public
// @Summary UserInfo
// @Produce application/json
// @Success 200 {string} echox.Response{data=models.UserInfo} "ok"
// @failure 400 {string} echox.Response "bad request"
// @failure 500 {string} echox.Response "internal error"
// @Router /api/v1/publics/user [get]
func (c PublicController) UserInfo(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)

	userInfo, err := c.userService.GetUserInfo(claims.ID)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: userInfo}.JSON(ctx)
}

// MenuTree
// @Tags Public
// @Summary UserMenuTree
// @Produce application/json
// @Success 200 {string} echox.Response{data=models.MenuTrees} "ok"
// @failure 400 {string} echox.Response "bad request"
// @failure 500 {string} echox.Response "internal error"
// @Router /api/v1/publics/user/menutree [get]
func (c PublicController) MenuTree(ctx echo.Context) error {
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)

	menuTrees, err := c.userService.GetUserMenuTrees(claims.ID)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: menuTrees}.JSON(ctx)
}

// UserLogin
// @Tags Public
// @Summary UserLogin
// @Produce application/json
// @Param data body dto.Login true "Login"
// @Success 200 {string} echox.Response "ok"
// @failure 400 {string} echox.Response "bad request"
// @failure 500 {string} echox.Response "internal error"
// @Router /api/v1/publics/user/login [post]
func (c PublicController) UserLogin(ctx echo.Context) error {
	login := new(dto.Login)

	if err := ctx.Bind(login); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	if c.config.Auth.Captcha.Enable {
		if login.CaptchaID == "" || login.CaptchaCode == "" {
			return echox.Response{Code: http.StatusBadRequest, Message: errors.CaptchaAnswerCodeEmpty}.JSON(ctx)
		}
		if !c.captcha.Verify(login.CaptchaID, login.CaptchaCode, false) {
			return echox.Response{Code: http.StatusBadRequest, Message: errors.CaptchaAnswerCodeNoMatch}.JSON(ctx)
		}
	}

	user, err := c.userService.Verify(login.Username, login.Password)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	token, err := c.authService.GenerateToken(user)
	if err != nil {
		return echox.Response{Code: http.StatusInternalServerError, Message: errors.AuthTokenGenerateFail}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"token": token}}.JSON(ctx)
}

// UserLogout
// @Tags Public
// @Summary UserLogout
// @Produce application/json
// @Success 200 {string} echox.Response "success"
// @Router /api/v1/publics/user/logout [post]
func (c PublicController) UserLogout(ctx echo.Context) error {
	claims, ok := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	if ok {
		_ = c.authService.DestroyToken(claims.Username)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// NewPublicController creates new public controller
func NewPublicController(
	userService services.UserService,
	authService services.AuthService,
	captcha lib.Captcha,
	logger lib.Logger,
	config lib.Config,
) PublicController {
	return PublicController{
		userService: userService,
		authService: authService,
		captcha:     captcha,
		logger:      logger,
		config:      config,
	}
}
