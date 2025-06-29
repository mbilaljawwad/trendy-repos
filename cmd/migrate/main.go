package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mbilaljawwad/trendy-repos/internal/config"
	"github.com/mbilaljawwad/trendy-repos/internal/datastore"
	"github.com/mbilaljawwad/trendy-repos/internal/migration"
)

func main() {
	var (
		action = flag.String("action", "up", "Migration action: up, down, status, create")
		name   = flag.String("name", "", "Migration name (required for create action)")
	)
	flag.Parse()

	// Initialize configuration
	config.InitConfig()

	// Get database connection (without running migrations automatically)
	ctx := context.Background()
	db := datastore.NewDataStoreWithoutMigrations(ctx)
	defer db.Close()

	// Create migrator
	migrator := migration.NewMigrator(db, "migrations")

	switch *action {
	case "up":
		if err := migrator.Up(); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	case "down":
		if err := migrator.Down(); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
	case "status":
		if err := migrator.Status(); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
	case "create":
		if *name == "" {
			fmt.Println("Migration name is required for create action")
			fmt.Println("Usage: go run cmd/migrate/main.go -action=create -name=migration_name")
			os.Exit(1)
		}
		if err := migrator.CreateMigrationFile(*name); err != nil {
			log.Fatalf("Failed to create migration file: %v", err)
		}
	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: up, down, status, create")
		os.Exit(1)
	}
}
