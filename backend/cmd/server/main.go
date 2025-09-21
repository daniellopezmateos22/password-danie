// main: composition root. Abre SQLite, cablea repos/usecases y registra rutas/readyz.
package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	api "password-danie/internal/http"
	"password-danie/internal/repository"
	sqliteRepo "password-danie/internal/repository/sqlite"
	"password-danie/internal/usecase"
	"password-danie/pkg/db"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("SQLITE_DSN")
	if dsn == "" {
		dsn = "data/app.db"
	}

	sqlDB, err := db.OpenSQLite(dsn)
	if err != nil {
		log.Fatalf("open sqlite: %v", err)
	}
	defer sqlDB.Close()

	// repos
	var (
		userRepo   repository.UserRepo   = sqliteRepo.NewUserSQLite(sqlDB)
		secretRepo repository.SecretRepo = sqliteRepo.NewSecretSQLite(sqlDB)
	)

	// usecases
	authUC := usecase.NewAuth(userRepo)
	vaultUC := usecase.NewVault(secretRepo)

	// http
	r := gin.Default()
	ready := func() error { return sqlDB.Ping() }

	api.RegisterRoutes(r, authUC, vaultUC, ready)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
