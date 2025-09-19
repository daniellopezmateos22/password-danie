package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"


	"github.com/daniellopezmateos22/password-danie/internal/models"
)

func Connect() *gorm.DB {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN no configurado")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("no puedo abrir Postgres: %v", err)
	}


	if err := db.AutoMigrate(
		&models.User{},
		&models.VaultItem{},
		&models.ResetToken{},
	); err != nil {
		log.Fatalf("migraciones: %v", err)
	}

	return db
}
