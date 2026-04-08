package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/theretech/retech-core-api/internal/domain"
	"github.com/theretech/retech-core-api/internal/interfaces/http/response"
)

const headerTenantID = "X-Tenant-ID"

func RequireTenant(repo domain.TenantRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader(headerTenantID)
		if raw == "" {
			response.Error(c, http.StatusBadRequest, "TENANT_REQUIRED", "header X-Tenant-ID is required")
			c.Abort()
			return
		}
		tid, err := uuid.Parse(raw)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "TENANT_INVALID", "X-Tenant-ID must be a valid UUID")
			c.Abort()
			return
		}
		tenant, err := repo.GetByID(c.Request.Context(), tid)
		if err != nil {
			if err == domain.ErrNotFound {
				response.Error(c, http.StatusNotFound, "TENANT_NOT_FOUND", "tenant not found")
				c.Abort()
				return
			}
			logAny, _ := c.Get(ContextKeyLogger)
			if log, ok := logAny.(*slog.Logger); ok {
				log.Error("tenant lookup failed", slog.Any("error", err))
			}
			response.Error(c, http.StatusInternalServerError, "INTERNAL", "internal error")
			c.Abort()
			return
		}
		if !tenant.Active {
			response.Error(c, http.StatusForbidden, "TENANT_INACTIVE", "tenant is inactive")
			c.Abort()
			return
		}
		c.Set(ContextKeyTenant, tenant)
		c.Next()
	}
}
