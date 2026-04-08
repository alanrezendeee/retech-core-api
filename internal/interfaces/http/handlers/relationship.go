package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/theretech/retech-core-api/internal/application/dto"
	"github.com/theretech/retech-core-api/internal/application/usecase"
	"github.com/theretech/retech-core-api/internal/interfaces/http/middleware"
	"github.com/theretech/retech-core-api/internal/interfaces/http/response"
)

type Relationship struct {
	UC *usecase.Relationship
}

func (h *Relationship) LinkApplication(c *gin.Context) {
	pathID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "tenant id must be a valid UUID")
		return
	}
	t, ok := middleware.CurrentTenant(c)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "INTERNAL", "tenant context missing")
		return
	}
	var req dto.LinkApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}
	if err := h.UC.LinkApplication(c.Request.Context(), pathID, t.ID, req.ApplicationID); err != nil {
		logUsecaseError(c, err)
		if response.FromDomainError(c, err) {
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL", "internal error")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Relationship) LinkUser(c *gin.Context) {
	pathID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "tenant id must be a valid UUID")
		return
	}
	t, ok := middleware.CurrentTenant(c)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "INTERNAL", "tenant context missing")
		return
	}
	var req dto.LinkUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}
	if err := h.UC.LinkUser(c.Request.Context(), pathID, t.ID, req.UserID); err != nil {
		logUsecaseError(c, err)
		if response.FromDomainError(c, err) {
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL", "internal error")
		return
	}
	c.Status(http.StatusNoContent)
}
