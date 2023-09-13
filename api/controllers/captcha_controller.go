package controllers

import (
	"github.com/labstack/echo/v4"
	"manuel71sj/go-api-template/errors"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models/dto"
	"manuel71sj/go-api-template/pkg/echox"
	"net/http"
)

type CaptchaController struct {
	captcha lib.Captcha
}

// GetCaptcha
// @Tags Public
// @Summary GetCaptcha
// @Produce application/json
// @Success 200 {string} echox.Response "ok"
// @failure 400 {string} echox.Response "bad request"
// @failure 500 {string} echox.Response "internal error"
// @Router /api/v1/publics/captcha [get]
func (c CaptchaController) GetCaptcha(ctx echo.Context) error {
	id, b64s, err := c.captcha.Generate()
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK, Data: echo.Map{"id": id, "blob": b64s}}.JSON(ctx)
}

// VerifyCaptcha
// @Tags Public
// @Summary VerifyCaptcha
// @Produce application/json
// @Param data body dto.CaptchaVerify true "CaptchaVerify"
// @Success 200 {string} echox.Response "ok"
// @failure 400 {string} echox.Response "bad request"
// @Router /api/v1/publics/captcha/verify [post]
func (c CaptchaController) VerifyCaptcha(ctx echo.Context) error {
	verify := new(dto.CaptchaVerify)
	if err := ctx.Bind(verify); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err}.JSON(ctx)
	}

	ok := c.captcha.Verify(verify.Id, verify.Code, false)
	if !ok {
		return echox.Response{Code: http.StatusBadRequest, Message: errors.CaptchaAnswerCodeNoMatch}.JSON(ctx)
	}

	return echox.Response{Code: http.StatusOK}.JSON(ctx)

}

// NewCaptchaController creates new captcha controller
func NewCaptchaController(captcha lib.Captcha) CaptchaController {
	return CaptchaController{captcha: captcha}
}
