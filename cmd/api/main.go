package main

import (
	"log"
	"os"
	"strconv"

	"github.com/Moshii-Xz/SeminarioIII/internal/platform/database"
	"github.com/Moshii-Xz/SeminarioIII/internal/server"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	db := database.New(database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "1234"),
		DBName:   getEnv("DB_NAME", "vencitrack"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
		TimeZone: getEnv("DB_TIMEZONE", "America/Bogota"),
	})

	//migrationsPath := filepath.Join(".", "migrations")
	//if err := database.RunMigrations(db, migrationsPath); err != nil {
	//	log.Fatalf("failed to run migrations: %v", err)
	//}

	servCfg := server.DefConfig()

	srv := server.New(db, servCfg)

	if err := srv.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
