package pg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"polley/models"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// PollController implementation for postgres.
type PollController struct {
	pool *pgxpool.Pool
}

// NewPollController constructs new Poll.
func NewPollController(pool *pgxpool.Pool) *PollController {
	return &PollController{
		pool,
	}
}

// Create creates a poll.
func (p *PollController) Create(poll *models.Poll) error {
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
		log.Printf("pg.PollController.Create: couldn't create %v: %v", poll, err)
		return err
	}

	for id, choice := range poll.Choices {
		_, err = tx.Exec(ctx, fmt.Sprintf(`INSERT INTO %s (id, poll_id, text, votes) VALUES ($1, $2, $3, $4)`, choicesTableName), id, pollID, choice.Text, choice.Votes)
		if err != nil {
			log.Printf("pg.PollController.Create: couldn't create %v: %v", poll, err)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("pg.PollController.Create: couldn't commit %v: %v", poll, err)
		return err
	}
	log.Printf("pg.PollController.Create: successfully created %v", poll)
	return nil
}

// Read reads a poll.
func (p *PollController) Read(uuid string) (*models.Poll, error) {
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
		log.Printf("pg.PollController.Read: couldn't read poll with uuid=%v: %v", uuid, err)
		return nil, err
	}

	choicesQuery, err := tx.Query(ctx, fmt.Sprintf(`SELECT text, votes FROM %s WHERE poll_id=(SELECT id FROM polls WHERE uuid=$1) ORDER BY id;`, choicesTableName), uuid)
	if err != nil {
		log.Printf("pg.PollController.Read: couldn't read choices for poll with uuid=%v: %v", uuid, err)
		return nil, err
	}
	for choicesQuery.Next() {
		choice := models.Choice{}
		err := choicesQuery.Scan(&choice.Text, &choice.Votes)
		if err != nil {
			log.Printf("pg.PollController.Read: couldn't read choice for poll with uuid=%v: %v", uuid, err)
			return nil, err
		}
		res.Choices = append(res.Choices, choice)
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("pg.PollController.Read: couldn't commit %v: %v", uuid, err)
		return nil, err
	}

	return res, nil
}

// Increment increments a choice for a poll with given uuid.
func (p *PollController) Increment(uuid string, choiceTexts []string) error {
	ctx := context.Background()

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, choiceText := range choiceTexts {
		_, err := tx.Exec(ctx, fmt.Sprintf(`UPDATE %s SET votes=votes+1 WHERE poll_id=(SELECT id FROM polls WHERE uuid=$1) AND text=$2`, choicesTableName), uuid, choiceText)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

// AddIPForPoll adds ip for poll.
func (p *PollController) AddIPForPoll(uuid string, ip string) error {
	ctx := context.Background()
	_, err := p.pool.Exec(ctx, `INSERT INTO ips (poll_id, ip) VALUES ((SELECT id FROM polls WHERE uuid=$1), $2) ON conflict DO nothing;`, uuid, ip)
	return err
}

// IsVoteAllowedForIP checks if vote is allowed for IP.
func (p *PollController) IsVoteAllowedForIP(uuid string, ip string) bool {
	ctx := context.Background()
	rows, err := p.pool.Query(ctx, `SELECT * FROM ips WHERE poll_id=(SELECT id FROM polls WHERE uuid=$1) AND ip=$2`, uuid, ip)
	if err != nil {
		return true
	}
	defer rows.Close()
	return !rows.Next()
}

// ErrInvalidArguments is thrown from function if invalid arguments are passed.
var ErrInvalidArguments = errors.New("invalid arguments")

// GetNPollsUUIDs gets n polls with limit and offset.
func (p *PollController) GetNPollsUUIDs(pageSize int, page int) ([]string, error) {
	if page < 0 || pageSize < 0 {
		return nil, ErrInvalidArguments
	}
	offset := page * pageSize
	rows, err := p.pool.Query(context.Background(), `SELECT uuid FROM polls ORDER BY id DESC LIMIT $1 OFFSET $2`, pageSize, offset)
	if err != nil {
		return nil, err
	}
	res := []string{}
	for rows.Next() {
		var uuid string
		err := rows.Scan(&uuid)
		if err != nil {
			log.Printf("GetNPollsUUIDs: error while scan: %v", err)
			return nil, err
		}
		res = append(res, uuid)
	}
	return res, nil
}

// GetTotal gets total entries.
func (p *PollController) GetTotal() int {
	row := p.pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM polls`)
	total := 0
	row.Scan(&total)
	return total
}
