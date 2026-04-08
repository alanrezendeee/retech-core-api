package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/theretech/retech-core-api/internal/domain"
)

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type envelope struct {
	Data  any        `json:"data,omitempty"`
	Error *errorBody `json:"error,omitempty"`
}

func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, envelope{Data: data})
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, envelope{Error: &errorBody{Code: code, Message: message}})
}

func FromDomainError(c *gin.Context, err error) bool {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		Error(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return true
	case errors.Is(err, domain.ErrConflict):
		Error(c, http.StatusConflict, "CONFLICT", err.Error())
		return true
	case errors.Is(err, domain.ErrInvalidInput):
		Error(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return true
	case errors.Is(err, domain.ErrForbidden):
		Error(c, http.StatusForbidden, "FORBIDDEN", err.Error())
		return true
	case errors.Is(err, domain.ErrTenantInactive):
		Error(c, http.StatusForbidden, "TENANT_INACTIVE", err.Error())
		return true
	case errors.Is(err, domain.ErrApplicationInactive):
		Error(c, http.StatusUnprocessableEntity, "APPLICATION_INACTIVE", err.Error())
		return true
	case errors.Is(err, domain.ErrNotLinked):
		Error(c, http.StatusForbidden, "NOT_LINKED", err.Error())
		return true
	default:
		return false
	}
}
