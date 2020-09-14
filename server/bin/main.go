package main

import (
	"log"
	"net/http"
	"polley/controllers/pg"
	"polley/server"
	"polley/server/storage.go"
)

var postgresURL = "postgres://postgres:postgrespass@localhost:5433/polley"

func main() {
	pool, err := pg.InitPool(postgresURL)
	if err != nil {
		log.Printf("pg.InitPool: %v", err)
	}
	pg.ApplyDefaultMigrations(pool)

	pgController := pg.NewPollController(pool)
	storage := storage.NewStorage(pgController, pgController)

	server := server.New(storage)

	log.Fatal(http.ListenAndServe(":5000", server))
}
