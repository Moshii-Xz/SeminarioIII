package main

import (
	"log"

	"github.com/mordmora/expirapp/internal/platform/database"
	"github.com/mordmora/expirapp/internal/server"
)

func main() {

	db := database.New(database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "1234",
		DBName:   "expirapp",
		SSLMode:  "disable",
		TimeZone: "America/Bogota",
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
