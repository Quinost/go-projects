package main

import (
	"database/sql"
	"flag"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"warehouse/internal/cfg"

	_ "github.com/lib/pq"
)

const (
	createMigrationTableSql = `
		CREATE TABLE IF NOT EXISTS schema_migrations (
        filename TEXT PRIMARY KEY,
        applied_at TIMESTAMP DEFAULT now()
    )`
	getAppliedMigrationSql = `SELECT filename FROM schema_migrations`
)

func main() {
	cfg, err := cfg.Load("../config.yaml")

	if err != nil {
		log.Fatalf("Error while reading config: %v", err)
	}

	flag.StringVar(&cfg.Database.ConnectionString, "db", cfg.Database.ConnectionString, "DB connection string")
	flag.BoolVar(&cfg.Database.Seed, "seed", cfg.Database.Seed, "Seed database")
	flag.Parse()

	createDbIfNotExist(cfg)

	db, err := sql.Open(cfg.Database.Driver, cfg.Database.ConnectionString)
	if err != nil {
		log.Fatalf("Error while connecting to db: %v", err)
	}

	defer db.Close()

	createMigrationTableIfNotExist(db)
	appliedMigrations := loadAppliedMigrations(db)

	populateDatabase(db, appliedMigrations, "migrations")
	if cfg.Database.Seed {
		populateDatabase(db, appliedMigrations, "seeds")
	}

}

func populateDatabase(db *sql.DB, appliedMigrations map[string]bool, folder string) {
	files, err := os.ReadDir(folder)
	if err != nil {
		log.Fatalf("Error reading %s error: %s", folder, err)
	}

	var migrationFileNames []string
	for _, entry := range files {
		if strings.HasSuffix(entry.Name(), ".sql") {
			migrationFileNames = append(migrationFileNames, entry.Name())
		}
	}

	for _, fileName := range migrationFileNames {
		path := filepath.Join(folder, fileName)
		if appliedMigrations[path] {
			continue
		}

		fileBytes, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Error reading file %v", err)
		}

		log.Printf("Pushing sql file: %s \n", path)

		_, err = db.Exec(string(fileBytes))
		if err != nil {
			log.Fatalf("Error while executing file %s error %v", path, err)
		}

		_, err = db.Exec(`INSERT INTO schema_migrations (filename) VALUES ($1)`, path)
		if err != nil {
			log.Fatalf("Failed to record migration %s: %v", path, err)
		}
	}
}

func createDbIfNotExist(cfg *cfg.Config) {
	dbName, postgres_cs := switchDatabase(cfg.Database.ConnectionString, cfg.Database.Driver)
	db, err := sql.Open(cfg.Database.Driver, postgres_cs)
	if err != nil {
		log.Fatalf("Error while connecting to db: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + dbName)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		log.Fatalf("Error creating database %s error: %v", dbName, err)
	}
}

func createMigrationTableIfNotExist(db *sql.DB) {
	_, err := db.Exec(createMigrationTableSql)
	if err != nil {
		log.Fatalf("Error while createing migration table: %v", err)
	}
}

func loadAppliedMigrations(db *sql.DB) map[string]bool {
	rows, err := db.Query(getAppliedMigrationSql)
	if err != nil {
		log.Fatalf("Error while loading applied migrations: %v", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatalf("Error while loading applied migrations: %v", err)
		}
		applied[name] = true
	}
	return applied
}

func switchDatabase(connStr string, driverName string) (originalDB string, adminConnStr string) {
	u, err := url.Parse(connStr)
	if err != nil {
		return "", ""
	}

	originalDB = strings.TrimPrefix(u.Path, "/")
	u.Path = "/" + driverName
	adminConnStr = u.String()
	return originalDB, adminConnStr
}
