package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	dbStatus := "up"

	// Check DB connection
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "down"
	}

	status := "ok"
	if dbStatus == "down" {
		status = "degraded"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"services": gin.H{
			"database": dbStatus,
		},
	})
}
