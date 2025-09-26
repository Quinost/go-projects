package main

import (
	"database/sql"
	"log"
	"warehouse/api"
	"warehouse/internal/cfg"
	"warehouse/internal/repositories"
	"warehouse/internal/services"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := cfg.Load("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config %v", err)
	}

	db, _ := sql.Open(cfg.Database.Driver, cfg.Database.ConnectionString)
	repo := repositories.InitializeRepositories(db)
	services := services.InitializeServices(cfg, repo)
	server := api.NewServer(cfg, services)

	if err := server.Run(); err != nil {
		log.Fatalf("Fatal error %v", err)
	}
}
