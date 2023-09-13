package dto

type CaptchaVerify struct {
	Id   string `json:"id" binding:"required"`
	Code string `json:"code" binding:"required"`
}
