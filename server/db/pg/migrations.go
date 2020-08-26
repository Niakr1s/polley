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
	fmt.Sprintf(`ALTER TABLE IF exists %s
		ADD COLUMN if not exists "filter" VARCHAR(10) default '';`, pollsTableName),
	fmt.Sprintf(`ALTER TABLE IF exists %s
		ADD CHECK(filter = 'ip' OR filter = 'cookie' OR filter = '')`, pollsTableName),
}

func applyMigrations(pool *pgxpool.Pool, migrations []string) {
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback(ctx)

	for i, m := range migrations {
		_, err := tx.Exec(ctx, m)
		if err != nil {
			log.Printf("pg.applyMigrations: exec#%v err: %v\n", i, err)
		} else {
			log.Printf("pg.applyMigrations: exec#%v success\n", i)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		panic(err)
	}
}
