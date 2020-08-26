package pg

import (
	"context"
	"fmt"
	"log"
	"polley/models"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// PollDB implementation for postgres.
type PollDB struct {
	pool *pgxpool.Pool
}

// NewPollDB constructs new Poll.
func NewPollDB(pool *pgxpool.Pool) *PollDB {
	return &PollDB{
		pool,
	}
}

// Create creates a poll in db.
func (p *PollDB) Create(poll *models.Poll) error {
	ctx := context.Background()
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx,
		fmt.Sprintf(`INSERT INTO %s (uuid, created_at, expires_at, allowMultiple, name, filter) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, pollsTableName),
		poll.UUID, time.Now().UTC().Round(time.Second), poll.Expires, poll.AllowMultiple, poll.Name, poll.Filter)
	var pollID int
	err = row.Scan(&pollID)
	if err != nil {
		log.Printf("pg.PollDB.Create: couldn't create %v: %v", poll, err)
		return err
	}

	for id, choice := range poll.Choices {
		_, err = tx.Exec(ctx, fmt.Sprintf(`INSERT INTO %s (id, poll_id, text, votes) VALUES ($1, $2, $3, $4)`, choicesTableName), id, pollID, choice.Text, choice.Votes)
		if err != nil {
			log.Printf("pg.PollDB.Create: couldn't create %v: %v", poll, err)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("pg.PollDB.Create: couldn't commit %v: %v", poll, err)
		return err
	}
	log.Printf("pg.PollDB.Create: successfully created %v", poll)
	return nil
}

// Read reads a poll from db.
func (p *PollDB) Read(uuid string) (*models.Poll, error) {
	ctx := context.Background()
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	res := models.NewEmptyPoll()

	pollQuery := tx.QueryRow(ctx, fmt.Sprintf(`SELECT uuid, expires_at, allowMultiple, name, filter FROM %s WHERE uuid=$1;`, pollsTableName), uuid)
	err = pollQuery.Scan(&res.UUID, &res.Expires, &res.AllowMultiple, &res.Name, &res.Filter)
	if err != nil {
		log.Printf("pg.PollDB.Read: couldn't read poll with uuid=%v: %v", uuid, err)
		return nil, err
	}

	choicesQuery, err := tx.Query(ctx, fmt.Sprintf(`SELECT text, votes FROM %s WHERE poll_id=(SELECT id FROM polls WHERE uuid=$1) ORDER BY id;`, choicesTableName), uuid)
	if err != nil {
		log.Printf("pg.PollDB.Read: couldn't read choices for poll with uuid=%v: %v", uuid, err)
		return nil, err
	}
	for choicesQuery.Next() {
		choice := models.Choice{}
		err := choicesQuery.Scan(&choice.Text, &choice.Votes)
		if err != nil {
			log.Printf("pg.PollDB.Read: couldn't read choice for poll with uuid=%v: %v", uuid, err)
			return nil, err
		}
		res.Choices = append(res.Choices, choice)
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("pg.PollDB.Read: couldn't commit %v: %v", uuid, err)
		return nil, err
	}

	return res, nil
}

// Increment increments a choice for a poll with given uuid.
func (p *PollDB) Increment(uuid string, choiceText string) error {
	ctx := context.Background()
	p.pool.Exec(ctx, fmt.Sprintf(`UPDATE %s	SET votes=votes+1 WHERE poll_id=(SELECT id FROM polls WHERE uuid=$1) AND text=$2`, choicesTableName), uuid, choiceText)

	return nil
}
