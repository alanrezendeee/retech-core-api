package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/theretech/retech-core-api/internal/application/dto"
	"github.com/theretech/retech-core-api/internal/application/usecase"
	"github.com/theretech/retech-core-api/internal/interfaces/http/middleware"
	"github.com/theretech/retech-core-api/internal/interfaces/http/response"
)

type Tenant struct {
	UC *usecase.Tenant
}

func (h *Tenant) Create(c *gin.Context) {
	var req dto.CreateTenantRequest
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

func (h *Tenant) ListCurrent(c *gin.Context) {
	t, ok := middleware.CurrentTenant(c)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "INTERNAL", "tenant context missing")
		return
	}
	out, err := h.UC.GetByID(c.Request.Context(), t.ID)
	if err != nil {
		if response.FromDomainError(c, err) {
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL", "internal error")
		return
	}
	response.JSON(c, http.StatusOK, []dto.TenantResponse{*out})
}

func logUsecaseError(c *gin.Context, err error) {
	logAny, ok := c.Get(middleware.ContextKeyLogger)
	if !ok {
		return
	}
	log, ok := logAny.(*slog.Logger)
	if !ok {
		return
	}
	log.Error("use case error", slog.Any("error", err))
}
