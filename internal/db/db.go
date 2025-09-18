package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/daniellopezmateos22/password-danie/internal/models"
)

func ConnectAndMigrate() *gorm.DB {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("no puedo abrir Postgres:", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.VaultItem{}, &models.ResetToken{}); err != nil {
		log.Fatal("no puedo migrar:", err)
	}
	return db
}
