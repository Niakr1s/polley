package db

import (
	"polley/models"
	"strings"
	"sync"
)

// PollDB is interface for storing PollDB.
type PollDB interface {
	Create(poll *models.Poll) error
	Read(uuid string) (*models.Poll, error)
	Increment(uuid string, choiceText string) error
	GetNPollsUUIDs(pageSize int, page int) ([]string, error)
}

// GetNPolls is a helper function, that combines PollDB.GetNPollsUUIDs and PollDB.Read.
func GetNPolls(pollDB PollDB, pageSize int, page int) ([]*models.Poll, error) {
	uuids, err := pollDB.GetNPollsUUIDs(pageSize, page)
	if err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(uuids))
	errors := []error{}
	errorCh := make(chan error)
	go func() {
		for err := range errorCh {
			errors = append(errors, err)
		}
	}()
	res := make([]*models.Poll, len(uuids))
	for i, uuid := range uuids {
		go func(idx int, uuid string) {
			defer wg.Done()
			poll, err := pollDB.Read(uuid)
			if err != nil {
				errorCh <- err
			}
			res[idx] = poll
		}(i, uuid)
	}
	wg.Wait()
	close(errorCh)
	if len(errors) != 0 {
		return nil, combineErrors(errors)
	}
	return res, nil
}

func combineErrors(errors []error) error {
	errorsStr := make([]string, len(errors))
	for i, err := range errors {
		errorsStr[i] = err.Error()
	}
	return CombinedError{errorsStr}
}

// CombinedError is a combined error.
type CombinedError struct {
	errors []string
}

func (e CombinedError) Error() string {
	return strings.Join(e.errors, "; ")
}
