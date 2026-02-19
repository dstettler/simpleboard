package main

import (
	"log"
	"net/http"

	"simpleboard/api"
	"simpleboard/pkg/config"
	"simpleboard/pkg/db"
)

func main() {
	// load env config
	cfg := config.Load()

	// connect the db
	db.Connect(cfg)

	// register API routes
	router := api.RegisterRoutes()

	log.Printf("Server running on %s\n", cfg.ServerAddress)

	if err := http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
