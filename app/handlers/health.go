package handlers

import (
	"context"
	"net/http"
	"strconv"

	"clienthub_go/config"
	// "tgithub.com/antoniomarfa/tools/api/middlewares"

	"github.com/gin-gonic/gin"
)

// SetHealthRoutes creates health routes
func SetHealthRoutes(ctx context.Context, cfg config.Config, r *gin.Engine) {
	r.GET("/health", healthCheck(cfg))
	r.HEAD("/health", healthCheck(cfg))
	r.GET("/", healthCheck(cfg))
	r.HEAD("/", healthCheck(cfg))
}

// @Summary Health Check
// @Description Runs a Health Check
// @Tags Health
// @Success 200 "OK"
// @Failure 500 {object} object
// @Failure 503 {object} object
// @Router /health [get]
// func healthCheck(ctx context.Context, cfg config.Config) http.Handler {
func healthCheck(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Agregar los headers con la información de configuración
		c.Header("Version", cfg.Version)
		c.Header("Environment", cfg.Environment)
		c.Header("Port", strconv.Itoa(cfg.Port))
		c.Header("Database", cfg.Database)
		c.Header("DSN", cfg.DSN)

		// Responder con un estado 200 OK y un cuerpo vacío (nil)
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}
}
