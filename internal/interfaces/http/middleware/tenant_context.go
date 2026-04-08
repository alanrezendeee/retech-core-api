package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/theretech/retech-core-api/internal/domain"
)

func CurrentTenant(c *gin.Context) (*domain.Tenant, bool) {
	v, ok := c.Get(ContextKeyTenant)
	if !ok {
		return nil, false
	}
	t, ok := v.(*domain.Tenant)
	return t, ok
}
