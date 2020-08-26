package pg

import (
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

var testPostgresURL = "postgres://postgres:postgrespass@localhost:5433/polley_test"

func makeTestMigrations() []string {
	res := []string{
		fmt.Sprintf("DROP TABLE IF EXISTS %s, %s", choicesTableName, pollsTableName),
	}
	res = append(res, defaultMigrations...)
	return res
}

func initTestPool() (*pgxpool.Pool, error) {
	pool, err := InitPool(testPostgresURL)
	return pool, err
}

func applyTestMigrations(pool *pgxpool.Pool) {
	applyMigrations(pool, makeTestMigrations())
}
