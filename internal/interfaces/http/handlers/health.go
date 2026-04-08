package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/theretech/retech-core-api/internal/infrastructure/persistence"
	"github.com/theretech/retech-core-api/internal/interfaces/http/response"
)

type Health struct {
	DB *gorm.DB
}

func (h *Health) Handle(c *gin.Context) {
	dbStatus := "up"
	if err := persistence.Ping(h.DB); err != nil {
		dbStatus = "down"
		response.JSON(c, http.StatusServiceUnavailable, gin.H{
			"status":   "degraded",
			"database": dbStatus,
		})
		return
	}
	response.JSON(c, http.StatusOK, gin.H{
		"status":   "ok",
		"database": dbStatus,
	})
}
