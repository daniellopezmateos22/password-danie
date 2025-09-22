// internal/http/router.go
// Registro de rutas HTTP con CORS (Codespaces/local), health/ready, auth, users/me y CRUD del vault.
package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"password-danie/internal/dto"
	"password-danie/internal/middleware"
	"password-danie/internal/usecase"
)

// RegisterRoutes registra endpoints públicos y protegidos sobre el *gin.Engine* recibido.
func RegisterRoutes(r *gin.Engine, authUC *usecase.Auth, vaultUC *usecase.Vault, readyCheck func() error) {
	// Evita warning de proxies y aplica CORS
	_ = r.SetTrustedProxies(nil)

	// CORS para local y Codespaces (*.app.github.dev)
	corsCfg := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	corsCfg.AllowOriginFunc = func(origin string) bool {
		// Local dev
		if origin == "http://localhost:5173" || origin == "https://localhost:5173" {
			return true
		}
		// GitHub Codespaces (cualquier subdominio)
		// Ej: https://<name>-5173.app.github.dev
		if strings.HasSuffix(origin, ".app.github.dev") {
			return true
		}
		return false
	}
	r.Use(cors.New(corsCfg))

	// Catch-all OPTIONS para preflight (evita 404 en preflight)
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// --- Health / Ready ---
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	r.GET("/readyz", func(c *gin.Context) {
		if readyCheck != nil && readyCheck() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"db": "down"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"db": "ok"})
	})

	api := r.Group("/api/v1")

	// --- Auth ---
	api.POST("/auth/register", func(c *gin.Context) {
		var req dto.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, err := authUC.Register(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, u)
	})

	api.POST("/auth/login", func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, u, err := authUC.Login(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"access_token": token, "user": u})
	})

	// --- Grupo protegido ---
	authGroup := api.Group("")
	authGroup.Use(middleware.AuthRequired())

	// users/me (lee claims del contexto que dejó el middleware)
	authGroup.GET("/users/me", func(c *gin.Context) {
		if claims, ok := c.Get("claims"); ok {
			c.JSON(http.StatusOK, gin.H{"claims": claims})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
	})

	// --- Vault CRUD ---
	v := authGroup.Group("/vault")

	v.GET("/entries", func(c *gin.Context) {
		uid := userIDFromClaims(c)
		q := c.Query("q")
		domain := c.Query("domain")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		items, total, err := vaultUC.List(uid, q, domain, limit, offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items, "total": total})
	})

	v.POST("/entries", func(c *gin.Context) {
		uid := userIDFromClaims(c)
		var req dto.CreateSecretRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, err := vaultUC.Create(uid, req.Username, req.PasswordPlain, req.URL, req.Notes, req.Icon, req.Title)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	})

	v.GET("/entries/:id", func(c *gin.Context) {
		uid := userIDFromClaims(c)
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		s, err := vaultUC.Get(uid, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if s == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, s)
	})

	v.PUT("/entries/:id", func(c *gin.Context) {
		uid := userIDFromClaims(c)
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		var req dto.UpdateSecretRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := vaultUC.Update(uid, id, req.Username, req.PasswordPlain, req.URL, req.Notes, req.Icon, req.Title); err != nil {
			if err.Error() == "not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	v.DELETE("/entries/:id", func(c *gin.Context) {
		uid := userIDFromClaims(c)
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		if err := vaultUC.Delete(uid, id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}

// userIDFromClaims obtiene el userID preferentemente del contexto (middleware) y si no, de los claims.
func userIDFromClaims(c *gin.Context) int64 {
	// 1) Preferir el valor que dejó el middleware
	if v, ok := c.Get(middleware.CtxUserID); ok {
		switch t := v.(type) {
		case int64:
			return t
		case int:
			return int64(t)
		case float64:
			return int64(t)
		}
	}
	// 2) Fallback: leer de los claims (por compatibilidad)
	claims, ok := c.Get("claims")
	if !ok {
		return 0
	}
	m, ok := claims.(map[string]any)
	if !ok {
		return 0
	}
	switch v := m["sub"].(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	}
	return 0
}
