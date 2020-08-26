package main

import (
	"log"
	"net/http"
	"polley/db/pg"
	"polley/server"
)

var postgresURL = "postgres://postgres:postgrespass@localhost:5433/polley"

func main() {
	pool, err := pg.InitPool(postgresURL)
	if err != nil {
		log.Printf("pg.InitPool: %v", err)
	}
	pg.ApplyDefaultMigrations(pool)

	pollDB := pg.NewPollDB(pool)

	server := server.New(pollDB, pollDB)

	log.Fatal(http.ListenAndServe(":5000", server))
}
