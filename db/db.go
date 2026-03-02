package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is a global variable holding the connection pool
var Pool *pgxpool.Pool

// InitDB connects to Supabase and pings the database
func InitDB() {
	// 1. Get the URL from the environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}

	// 2. Connect to the database pool
	var err error
	Pool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// 3. Pin the database to ensure connection works
	err = Pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database didn't respond to ping: %v\n", err)
	}

	log.Println("Successfully connnected to the Supabase database!")
}
