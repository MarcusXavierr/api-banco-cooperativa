package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/db"
	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/router"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	pool := createConnectionPool()
	defer pool.Close()

	queries := db.New(pool)
	router.Initialize(queries)
}

func createConnectionPool() *pgxpool.Pool {
	ctx := context.Background()

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, dbHost, dbPort, dbName)
	config := connectionString + "?pool_min_conns=4"
	conn, _ := pgxpool.ParseConfig(config)

	pool, err := pgxpool.NewWithConfig(ctx, conn)

	if err != nil {
		log.Fatalf("Error while connection to database: %v\n", err)
	}

	return pool
}
