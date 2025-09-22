// Handlers HTTP para el flujo de recuperación de contraseña (request + confirm).
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"password-danie/internal/usecase"
)

func RegisterResetRoutes(r *gin.Engine, resetUC *usecase.PasswordReset) {
	api := r.Group("/api/v1/auth/reset")

	api.POST("/request", func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required,email"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		token, err := resetUC.Request(req.Email)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// aquí se enviaría por email
		c.JSON(http.StatusOK, gin.H{"reset_token": token})
	})

	api.POST("/confirm", func(c *gin.Context) {
		var req struct {
			Token       string `json:"token" binding:"required"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := resetUC.Confirm(req.Token, req.NewPassword); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}
