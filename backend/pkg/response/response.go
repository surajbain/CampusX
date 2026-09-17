package response

import (
	"net/http"

	"github.com/campusx/api/pkg/apperror"
	"github.com/gin-gonic/gin"
)

type SuccessEnvelope struct {
	Success bool `json:"success"`
	Data    any  `json:"data,omitempty"`
	Meta    any  `json:"meta,omitempty"`
}

type ErrorEnvelope struct {
	Success bool      `json:"success"`
	Error   ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, SuccessEnvelope{Success: true, Data: data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, SuccessEnvelope{Success: true, Data: data})
}

func WithMeta(c *gin.Context, status int, data, meta any) {
	c.JSON(status, SuccessEnvelope{Success: true, Data: data, Meta: meta})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Fail converts any error to a standard error response.
// Internal error details are logged, never exposed.
func Fail(c *gin.Context, err error) {
	ae := apperror.As(err)
	if ae == nil {
		c.JSON(http.StatusInternalServerError, ErrorEnvelope{
			Success: false,
			Error:   ErrorBody{Code: string(apperror.CodeInternal), Message: "internal error"},
		})
		return
	}
	c.JSON(ae.HTTPStatus(), ErrorEnvelope{
		Success: false,
		Error: ErrorBody{
			Code:    string(ae.Code),
			Message: ae.Message,
			Details: ae.Details,
		},
	})
}
