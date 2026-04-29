package main

import (
	"log"
	"net/http"
	"time"

	"simpleboard/api"
	"simpleboard/internal/auth"
	"simpleboard/internal/timer"
	"simpleboard/pkg/config"
	"simpleboard/pkg/db"
)

func main() {
	// load env config
	cfg := config.Load()

	// initialize auth middleware
	auth.Init(cfg)

	// load timer defaults from config
	timer.Init(cfg)

	// connect the db
	db.Connect(cfg)

	// kick off background sweeper that ends games whose clocks ran out
	// even if no one is polling them
	timer.StartSweeper(db.DB, time.Duration(cfg.SweepIntervalSeconds)*time.Second)

	// register API routes
	router := api.RegisterRoutes()

	log.Printf("Server running on %s\n", cfg.ServerAddress)

	if err := http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
