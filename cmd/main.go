package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/daniellopezmateos22/password-danie/api"
	dbpkg "github.com/daniellopezmateos22/password-danie/internal/db"
)

func main() {
	var gdb *gorm.DB = dbpkg.ConnectAndMigrate()

	r := gin.Default()
	// CORS abierto para dev; en prod, restringe Origins
	r.Use(cors.Default())

	api.Routes(r, gdb)

	log.Println("listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
