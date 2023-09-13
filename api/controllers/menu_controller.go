package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"manuel71sj/go-api-template/api/services"
	"manuel71sj/go-api-template/constants"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models"
	"manuel71sj/go-api-template/models/dto"
	"manuel71sj/go-api-template/pkg/echox"
	"net/http"
)

type MenuController struct {
	menuService services.MenuService
	logger      lib.Logger
}

// Query
// @Tags Menu
// @Summary Query
// @Produce application/json
// @Param data query models.MenuQueryParam true "MenuQueryParam"
// @Success 200 {object} echox.Response{data=models.MenuQueryResult} "ok"
// @failure 400 {string} echox.Response "bad request"
// @failure 500 {string} echox.Response "internal error"
// @Router /api/v1/menus [get]
func (c MenuController) Query(ctx echo.Context) error {
	param := new(models.MenuQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := c.menuService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	if param.Tree {
		return echox.Response{Code: http.StatusOK, Data: qr.List.ToMenuTrees()}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// Get
// @Tags Menu
// @Summary Menu Get By ID
// @Produce application/json
// @Param id path int true "menu id"
// @Success 200 {object} echox.Response{data=models.Menu} "ok
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/menus/{id} [get]
func (c MenuController) Get(ctx echo.Context) error {
	menu, err := c.menuService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: menu}.JSON(ctx)
}

// Create
// @Tags Menu
// @Summary Menu Create
// @Produce application/json
// @Param data body models.Menu true "Menu"
// @Success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/menus [post]
func (c MenuController) Create(ctx echo.Context) error {
	menu := new(models.Menu)
	if err := ctx.Bind(menu); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	menu.CreatedBy = claims.Username

	id, err := c.menuService.WithTrx(trxHandle).Create(menu)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// Update
// @Tags Menu
// @Summary Menu Update By ID
// @Produce application/json
// @Param id path int true "menu id"
// @Param data body models.Menu true "Menu"
// @Success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/menus/{id} [put]
func (c MenuController) Update(ctx echo.Context) error {
	menu := new(models.Menu)
	if err := ctx.Bind(menu); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := c.menuService.WithTrx(trxHandle).Update(ctx.Param("id"), menu); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// Delete
// @Tags Menu
// @Summary Menu Delete By ID
// @Produce application/json
// @Param id path int true "menu id"
// @Success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/menus/{id} [delete]
func (c MenuController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := c.menuService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// Enable
// @Tags Menu
// @Summary Menu Enable By ID
// @Produce application/json
// @Param id path int true "menu id"
// @Success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/menus/{id}/enable [patch]
func (c MenuController) Enable(ctx echo.Context) error {
	if err := c.menuService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// Disable
// @Tags Menu
// @Summary Menu Disable By ID
// @Produce application/json
// @Param id path int true "menu id"
// @Success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/menus/{id}/disable [patch]
func (c MenuController) Disable(ctx echo.Context) error {
	if err := c.menuService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// GetActions
// @Tags Menu
// @Summary MenuActions Get By ID
// @Produce application/json
// @Param id path int true "menu id"
// @Success 200 {object} echox.Response{data=models.MenuActions} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/menus/{id}/actions [get]
func (c MenuController) GetActions(ctx echo.Context) error {
	actions, err := c.menuService.GetMenuActions(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: actions}.JSON(ctx)
}

// UpdateActions
// @Tags Menu
// @Summary MenuActions Update By ID
// @Produce application/json
// @Param id path int true "menu id"
// @Param data body models.MenuActions true "MenuActions"
// @Success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/menus/{id}/actions [put]
func (c MenuController) UpdateActions(ctx echo.Context) error {
	actions := make(models.MenuActions, 0)
	if err := ctx.Bind(&actions); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := c.menuService.WithTrx(trxHandle).UpdateActions(ctx.Param("id"), actions); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// NewMenuController creates new menu controller
func NewMenuController(
	logger lib.Logger,
	menuService services.MenuService,
) MenuController {
	return MenuController{
		logger:      logger,
		menuService: menuService,
	}
}
