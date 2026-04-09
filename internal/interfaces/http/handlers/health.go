package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/theretech/retech-core-api/internal/version"
)

type healthResponse struct {
	Service  string `json:"service"`
	Status   string `json:"status"`
	DataBase string `json:"dataBase"`
	Version  string `json:"version"`
}

type Health struct {
	DB *gorm.DB
}

func (h *Health) Handle(c *gin.Context) {
	dataBase := "up"
	status := "ok"
	code := http.StatusOK

	sqlDB, err := h.DB.DB()
	if err != nil {
		dataBase = "down"
		status = "degraded"
		code = http.StatusServiceUnavailable
	} else if err := sqlDB.PingContext(c.Request.Context()); err != nil {
		dataBase = "down"
		status = "degraded"
		code = http.StatusServiceUnavailable
	}

	c.JSON(code, healthResponse{
		Service:  version.Service,
		Status:   status,
		DataBase: dataBase,
		Version:  version.Version,
	})
}
