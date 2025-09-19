package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	mydb "github.com/daniellopezmateos22/password-danie/internal/db"
	"github.com/daniellopezmateos22/password-danie/api"
)

func main() {
	// Router base con logger y recovery
	r := gin.Default()

	// CORS global
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Catch-all para cualquier preflight OPTIONS (evita 404)
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Length,Content-Type,Authorization")
		c.Status(204)
	})

	// DB + rutas
	db := mydb.Connect()
	api.Routes(r, db)

	log.Println("listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
