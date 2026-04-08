package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/theretech/retech-core-api/internal/application/dto"
	"github.com/theretech/retech-core-api/internal/application/usecase"
	"github.com/theretech/retech-core-api/internal/interfaces/http/middleware"
	"github.com/theretech/retech-core-api/internal/interfaces/http/response"
)

type Application struct {
	UC *usecase.Application
}

func (h *Application) Create(c *gin.Context) {
	var req dto.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}
	out, err := h.UC.Create(c.Request.Context(), req)
	if err != nil {
		logUsecaseError(c, err)
		if response.FromDomainError(c, err) {
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL", "internal error")
		return
	}
	response.JSON(c, http.StatusCreated, out)
}

func (h *Application) ListLinked(c *gin.Context) {
	t, ok := middleware.CurrentTenant(c)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "INTERNAL", "tenant context missing")
		return
	}
	list, err := h.UC.ListForTenant(c.Request.Context(), t.ID)
	if err != nil {
		logUsecaseError(c, err)
		if response.FromDomainError(c, err) {
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL", "internal error")
		return
	}
	response.JSON(c, http.StatusOK, list)
}
