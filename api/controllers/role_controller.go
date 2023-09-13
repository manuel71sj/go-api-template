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

type RoleController struct {
	logger      lib.Logger
	roleService services.RoleService
}

// Query
// @Tags Role
// @Summary Role Query
// @Produce application/json
// @Param data query models.RoleQueryParam true "RoleQueryParam"
// @Success 200 {object} echox.Response{data=models.RoleQueryResult} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/roles [get]
func (c RoleController) Query(ctx echo.Context) error {
	param := new(models.RoleQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	qr, err := c.roleService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// GetAll
// @Tags Role
// @Summary Role Get All
// @Produce application/json
// @Param data query models.RoleQueryParam true "RoleQueryParam"
// @Success 200 {object} echox.Response{data=models.Roles} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/roles.all [get]
func (c RoleController) GetAll(ctx echo.Context) error {
	qr, err := c.roleService.Query(&models.RoleQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr.List}.JSON(ctx)
}

// Get
// @Tags Role
// @Summary Role Get By ID
// @Produce application/json
// @Param id path int true "role id"
// @Success 200 {object} echox.Response{data=models.Role} "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/roles/{id} [get]
func (c RoleController) Get(ctx echo.Context) error {
	role, err := c.roleService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: role}.JSON(ctx)
}

// Create
// @Tags Role
// @Summary Role Create
// @Produce application/json
// @Param data body models.Role true "Role"
// @Success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/roles [post]
func (c RoleController) Create(ctx echo.Context) error {
	role := new(models.Role)
	if err := ctx.Bind(role); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	claims, _ := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
	role.CreatedBy = claims.Username

	id, err := c.roleService.WithTrx(trxHandle).Create(role)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id}}.JSON(ctx)
}

// Update
// @Tags Role
// @Summary Role Update By ID
// @Produce application/json
// @Param id path int true "role id"
// @Param data body models.Role true "Role"
// @Success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/roles/{id} [put]
func (c RoleController) Update(ctx echo.Context) error {
	role := new(models.Role)
	if err := ctx.Bind(role); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := c.roleService.WithTrx(trxHandle).Update(ctx.Param("id"), role); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// Delete
// @Tags Role
// @Summary Role Delete By ID
// @Produce application/json
// @Param id path int true "role id"
// @Success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/roles/{id} [delete]
func (c RoleController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)
	if err := c.roleService.WithTrx(trxHandle).Delete(ctx.Param("id")); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// Enable
// @Tags Role
// @Summary Role Enable By ID
// @Produce application/json
// @Param id path int true "role id"
// @Success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/roles/{id}/enable [put]
func (c RoleController) Enable(ctx echo.Context) error {
	if err := c.roleService.UpdateStatus(ctx.Param("id"), 1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// Disable
// @Tags Role
// @Summary Role Disable By ID
// @Produce application/json
// @Param id path int true "role id"
// @Success 200 {object} echox.Response "ok"
// @failure 400 {object} echox.Response "bad request"
// @failure 500 {object} echox.Response "internal error"
// @Router /api/v1/roles/{id}/disable [put]
func (c RoleController) Disable(ctx echo.Context) error {
	if err := c.roleService.UpdateStatus(ctx.Param("id"), -1); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// NewRoleController creates a new role controller
func NewRoleController(
	logger lib.Logger,
	roleService services.RoleService,
) RoleController {
	return RoleController{
		logger:      logger,
		roleService: roleService,
	}
}
