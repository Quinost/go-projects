package main

import (
	"database/sql"
	"log"
	"warehouse/api"
	"warehouse/internal/config"
	"warehouse/internal/repositories"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config %v", err)
	}

	db, _ := sql.Open(cfg.Database.Driver, cfg.Database.ConnectionString)
	repo := repositories.InitializeRepositories(db)
	server := api.NewServer(repo)

	if err := server.Run(); err != nil {
		log.Fatal("Fatal error", err)
	}
}
