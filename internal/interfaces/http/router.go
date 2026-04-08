package httptransport

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/theretech/retech-core-api/internal/application/usecase"
	"github.com/theretech/retech-core-api/internal/domain"
	"github.com/theretech/retech-core-api/internal/interfaces/http/handlers"
	"github.com/theretech/retech-core-api/internal/interfaces/http/middleware"
)

type RouterDeps struct {
	DB              *gorm.DB
	Log             *slog.Logger
	TenantRepo      domain.TenantRepository
	TenantUC        *usecase.Tenant
	ApplicationUC   *usecase.Application
	RelationshipUC  *usecase.Relationship
	GinMode         string
}

func NewRouter(d RouterDeps) *gin.Engine {
	gin.SetMode(d.GinMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID(d.Log))
	r.Use(middleware.AccessLog())

	healthH := &handlers.Health{DB: d.DB}
	r.GET("/health", healthH.Handle)

	tenantH := &handlers.Tenant{UC: d.TenantUC}
	r.POST("/tenants", tenantH.Create)

	protected := r.Group("")
	protected.Use(middleware.RequireTenant(d.TenantRepo))
	{
		protected.GET("/tenants", tenantH.ListCurrent)

		appH := &handlers.Application{UC: d.ApplicationUC}
		protected.POST("/applications", appH.Create)
		protected.GET("/applications", appH.ListLinked)

		relH := &handlers.Relationship{UC: d.RelationshipUC}
		protected.POST("/tenants/:id/applications", relH.LinkApplication)
		protected.POST("/tenants/:id/users", relH.LinkUser)
	}

	return r
}
