package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yosa/ocr-golang-back/api"
	"github.com/yosa/ocr-golang-back/db"
	"github.com/yosa/ocr-golang-back/util"
)

func main() {
	// Load env vars using Viper
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Connect to PostgreSQL
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to DB: %v", err)
	}
	defer conn.Close()

	queries := db.New(conn)

	// Create server
	server, err := api.NewServer(config, queries)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}

	// Start server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
