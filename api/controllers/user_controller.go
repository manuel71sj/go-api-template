package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"manuel71sj/go-api-template/api/services"
	"manuel71sj/go-api-template/constants"
	"manuel71sj/go-api-template/errors"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models"
	"manuel71sj/go-api-template/models/dto"
	"manuel71sj/go-api-template/pkg/echox"
	"net/http"
	"strings"
)

type UserController struct {
	userService services.UserService
	logger      lib.Logger
}

// Query
// @Tags User
// @Summary User Query
// @Produce application/json
// @Param data query models.UserQueryParam true "UserQueryParam"
// @Success 200 {object} echox.Response{data=models.UserQueryResult} "ok"
// @Failure 400 {object} echox.Response "bad request"
// @Failure 500 {object} echox.Response "internal server error"
// @Router /api/v1/users [get]
func (c UserController) Query(ctx echo.Context) error {
	param := new(models.UserQueryParam)
	if err := ctx.Bind(param); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	if v := ctx.QueryParam("role_ids"); v != "" {
		param.RoleIDs = strings.Split(v, ",")
	}

	qr, err := c.userService.Query(param)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// Create
// @Tags User
// @Summary User Create
// @Produce application/json
// @Param data body models.User true "User"
// @Success 200 {object} echox.Response "ok"
// @Failure 400 {object} echox.Response "bad request"
// @Failure 500 {object} echox.Response "internal server error"
// @Router /api/v1/users [post]
func (c UserController) Create(ctx echo.Context) error {
	user := new(models.User)
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)

	if err := ctx.Bind(user); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	} else if user.Password == "" {
		return echox.Response{Code: http.StatusBadRequest, Message: errors.UserPasswordRequired}.JSON(ctx)
	}

	claims, _ := ctx.Get(constants.CurrentUser).(dto.JwtClaims)
	user.CreatedBy = claims.Username

	qr, err := c.userService.WithTrx(trxHandle).Create(user)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: qr}.JSON(ctx)
}

// Get
// @Tags User
// @Summary User Get By ID
// @Produce application/json
// @Param id path int true "user id"
// @Success 200 {object} echox.Response{data=models.User} "ok"
// @Failure 400 {object} echox.Response "bad request"
// @Failure 500 {object} echox.Response "internal server error"
// @Router /api/v1/users/{id} [get]
func (c UserController) Get(ctx echo.Context) error {
	user, err := c.userService.Get(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: user}.JSON(ctx)
}

// Update
// @Tags User
// @Summary User Update By ID
// @Produce application/json
// @Param id path int true "user id"
// @Param data body models.User true "User"
// @Success 200 {object} echox.Response "ok"
// @Failure 400 {object} echox.Response "bad request"
// @Failure 500 {object} echox.Response "internal server error"
// @Router /api/v1/users/{id} [put]
func (c UserController) Update(ctx echo.Context) error {
	user := new(models.User)
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)

	if err := ctx.Bind(user); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	err := c.userService.WithTrx(trxHandle).Update(ctx.Param("id"), user)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// Delete
// @Tags User
// @Summary User Delete By ID
// @Produce application/json
// @Param id path int true "user id"
// @Success 200 {object} echox.Response "ok"
// @Failure 400 {object} echox.Response "bad request"
// @Failure 500 {object} echox.Response "internal server error"
// @Router /api/v1/users/{id} [delete]
func (c UserController) Delete(ctx echo.Context) error {
	trxHandle := ctx.Get(constants.DBTransaction).(*gorm.DB)

	err := c.userService.WithTrx(trxHandle).Delete(ctx.Param("id"))
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// Enable
// @Tags User
// @Summary User Enable By ID
// @Produce application/json
// @Param id path int true "user id"
// @Success 200 {object} echox.Response "ok"
// @Failure 400 {object} echox.Response "bad request"
// @Failure 500 {object} echox.Response "internal server error"
// @Router /api/v1/users/{id}/enable [put]
func (c UserController) Enable(ctx echo.Context) error {
	err := c.userService.UpdateStatus(ctx.Param("id"), 1)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// Disable
// @Tags User
// @Summary User Disable By ID
// @Produce application/json
// @Param id path int true "user id"
// @Success 200 {object} echox.Response "ok"
// @Failure 400 {object} echox.Response "bad request"
// @Failure 500 {object} echox.Response "internal server error"
// @Router /api/v1/users/{id}/disable [put]
func (c UserController) Disable(ctx echo.Context) error {
	err := c.userService.UpdateStatus(ctx.Param("id"), -1)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)
}

// NewUserController creates new user controller
func NewUserController(
	userService services.UserService,
	logger lib.Logger,
) UserController {
	return UserController{
		userService: userService,
		logger:      logger,
	}
}
