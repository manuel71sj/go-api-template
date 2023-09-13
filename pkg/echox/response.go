package exhox

import (
	"github.com/labstack/echo/v4"
	"manuel71sj/go-api-template/errors"
	"net/http"
)

// Response in order to unify the returned response structure
type Response struct {
	Code    int         `json:"-"`
	Pretty  bool        `json:"-"`
	Data    interface{} `json:"data,omitempty"`
	Message interface{} `json:"message"`
}

// JSON sends a JSON response with status code
func (r *Response) JSON(ctx echo.Context) error {
	if r.Message == "" || r.Message == nil {
		r.Message = http.StatusText(r.Code)
	}

	if err, ok := r.Message.(error); ok {
		if errors.Is(err, errors.DatabaseInternalError) {
			r.Code = http.StatusInternalServerError
		}

		if errors.Is(err, errors.DatabaseRecordNotFound) {
			r.Code = http.StatusNotFound
		}

		r.Message = err.Error()
	}

	if r.Pretty {
		return ctx.JSONPretty(r.Code, r, "\t")
	}

	return ctx.JSON(r.Code, r)

}
