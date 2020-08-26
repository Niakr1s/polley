package pg

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var pollsTableName = "polls"
var choicesTableName = "choices"

var defaultMigrations []string = []string{
	fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id SERIAL PRIMARY KEY,
		uuid VARCHAR(36) NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL,
		expires_at TIMESTAMP NOT NULL
	);`, pollsTableName),
	fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INT NOT NULL,
		poll_id INT,
		text VARCHAR(100) NOT NULL,
		votes INT DEFAULT 0,
		FOREIGN KEY (poll_id) REFERENCES polls(id) ON DELETE CASCADE ON UPDATE CASCADE,
		UNIQUE (poll_id, id, text)
	);`, choicesTableName),
	fmt.Sprintf(`ALTER TABLE IF exists %s
		ADD COLUMN if not exists allowMultiple INT NOT NULL DEFAULT 1
		CHECK (allowMultiple >= 1);`, pollsTableName),
	fmt.Sprintf(`ALTER TABLE IF exists %s
		ADD COLUMN if not exists "name" VARCHAR(100);`, pollsTableName),
}

func applyMigrations(pool *pgxpool.Pool, migrations []string) {
	ctx := context.Background()

	for i, m := range migrations {
		_, err := pool.Exec(ctx, m)
		if err != nil {
			log.Printf("pg.applyMigrations: exec#%v err: %v\n", i, err)
		} else {
			log.Printf("pg.applyMigrations: exec#%v success\n", i)
		}
	}
}