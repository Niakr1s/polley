package controllers

import (
	"polley/models"
)

// PollController is interface for storing PollController.
type PollController interface {
	Create(poll *models.Poll) error
	Read(uuid string) (*models.Poll, error)
	Increment(uuid string, choiceTexts []string) error
	GetNPollsUUIDs(pageSize int, page int) ([]string, error)
	GetTotal() int
}
