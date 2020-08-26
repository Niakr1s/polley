package pg

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// InitPool inits new Pool.
func InitPool(url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}
	log.Printf("pg.InitPool: succesfully connected to database")

	return pool, nil
}

// ApplyDefaultMigrations creates nessesarry tables after init pool.
func ApplyDefaultMigrations(pool *pgxpool.Pool) {
	applyMigrations(pool, defaultMigrations)
}
