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

	db, err := sql.Open(cfg.Database.Driver, cfg.Database.ConnectionString)
	if err != nil {
		log.Fatalf("Failed opening connection to db %v", err)
	}

	repo := repo.InitializeRepositories(db)
	services := services.InitializeServices(cfg, repo)
	server := api.NewServer(cfg, services)

	errChan := make(chan error)

	go func() {
		if err := server.Run(); err != nil {
			errChan <-err
		}
	}()

	err = <-errChan

	log.Fatalf("Fatal error %v", err)
}
