package handlers

import (
	"context"
	"net/http"
	"strings"

	"clienthub_go/config"
	"clienthub_go/util"

	"github.com/gin-gonic/gin"
)

// SendCodeRequest representa el payload recibido para enviar el código de verificación
type SendCodeRequest struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// SetEmailRoutes registra las rutas para el envío de correos
func SetEmailRoutes(ctx context.Context, cfg config.Config, router *gin.Engine) {
	emailGroup := router.Group("api/v2027_01/")
	{
		// Ruta solicitada: POST /api/v2027_01/email/send-code
		emailGroup.POST("/email/send-code", sendCodeHandler(cfg))
		// Alias conveniente por si el frontend invoca /api/v2027_01/send-code
		emailGroup.POST("/send-code", sendCodeHandler(cfg))
	}
}

// @Summary Enviar código de verificación por email
// @Description Envía un correo electrónico con el código de verificación solicitado
// @Tags email
// @Accept json
// @Produce json
// @Param request body handlers.SendCodeRequest true "Email y código a enviar"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /api/v2027_01/email/send-code [post]
func sendCodeHandler(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SendCodeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Faltan datos obligatorios (email y code)",
			})
			return
		}

		req.Email = strings.TrimSpace(req.Email)
		req.Code = strings.TrimSpace(req.Code)

		if req.Email == "" || req.Code == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "El correo y el código no pueden estar vacíos",
			})
			return
		}

		// Enviar el correo usando la función en util/email.go
		if err := util.SendVerificationCodeEmail(cfg, req.Email, req.Code); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Error al enviar el correo: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Código de verificación enviado exitosamente",
		})
	}
}
