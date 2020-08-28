package db

import (
	"polley/models"
)

// PollDB is interface for storing PollDB.
type PollDB interface {
	Create(poll *models.Poll) error
	Read(uuid string) (*models.Poll, error)
	Increment(uuid string, choiceText string) error
	GetNPollsUUIDs(pageSize int, page int) ([]string, error)
	GetTotal() int
}
