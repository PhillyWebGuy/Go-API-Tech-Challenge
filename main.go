package main

import (
	"log"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/database"
	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/webserver"
)

func main() {
	// ctx := context.Background()
	if err := run( /*ctx*/ ); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func run( /*ctx context.Context*/ ) error {
	// Initialize the database
	db, err := database.Init()
	if err != nil {
		return err
	}

	// Pass the database connection to the web server
	webserver.NewServer(db)
	return nil
}
