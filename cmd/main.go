package main

import (
	"log"
	"os"
	"strconv"

	"github.com/junedAtPurplesky/purplegate-merchant-service/internal/db"
	"github.com/junedAtPurplesky/purplegate-merchant-service/internal/server"
)

func main() {
	// Get port from environment variable or use default
	portStr := os.Getenv("GRPC_PORT")
	if portStr == "" {
		portStr = "50051" // Default gRPC port
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	// Initialize database configuration
	dbConfig := db.NewConfig()
	log.Printf("Connecting to database: %s@%s:%s/%s", dbConfig.User, dbConfig.Host, dbConfig.Port, dbConfig.DBName)

	// Connect to database
	pool, err := db.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close(pool)

	// Run database migrations
	migrationsPath := "migrations"
	if err := db.RunMigrations(pool, migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create merchant repository
	merchantRepo := db.NewMerchantRepository(pool)

	// Create and start the merchant server
	merchantServer := server.NewMerchantServer(merchantRepo)

	log.Printf("Starting PurpleGate Merchant Service on port %d", port)
	if err := merchantServer.StartServer(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
