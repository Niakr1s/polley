package storage

import (
	"errors"
	"polley/models"
	"sync"
)

// NewMockStorage constructs storage with in-memory implementations.
func NewMockStorage() *Storage {
	return &Storage{
		Polls: &mockPollsController{},
		Ips:   &mockIpsController{},
	}
}

type mockPollsController struct {
	sync.RWMutex
	polls []*models.Poll
}

func newMockPollsController() *mockPollsController {
	return &mockPollsController{
		polls: make([]*models.Poll, 0),
	}
}

func (c *mockPollsController) Create(poll *models.Poll) error {
	c.Lock()
	defer c.Unlock()

	_, err := c.findPoll(poll.UUID)
	if err == nil {
		return errors.New("poll with such uuid already exists")
	}

	c.polls = append(c.polls, poll)
	return nil
}

func (c *mockPollsController) Read(uuid string) (*models.Poll, error) {
	c.RLock()
	defer c.RUnlock()

	return c.findPoll(uuid)
}
func (c *mockPollsController) Increment(uuid string, choiceTexts []string) error {
	c.Lock()
	defer c.Unlock()

	poll, err := c.findPoll(uuid)
	if err != nil {
		return err
	}

	// checking if all choiceTexts exists
	for _, choiceText := range choiceTexts {
		found := false
		for _, choice := range poll.Choices {
			if choice.Text == choiceText {
				found = true
			}
		}
		if !found {
			return errors.New("couldn't find one of choicetexts in poll")
		}
	}

	// incrementing each choice
	for _, choiceText := range choiceTexts {
		for i, choice := range poll.Choices {
			if choice.Text == choiceText {
				choice.Votes++
				poll.Choices[i] = choice
			}
		}
	}
	return nil
}

func (c *mockPollsController) GetNPollsUUIDs(pageSize int, page int) ([]string, error) {
	c.RLock()
	defer c.RUnlock()

	start := page * pageSize
	end := start + pageSize
	if end > len(c.polls) {
		end = len(c.polls)
	}

	if start < 0 || start > end {
		return nil, errors.New("index out of bounds")
	}

	polls := c.polls[start:end]
	result := []string{}
	for _, poll := range polls {
		result = append(result, poll.UUID)
	}
	return result, nil
}

func (c *mockPollsController) GetTotal() int {
	c.RLock()
	defer c.RUnlock()

	return len(c.polls)
}

func (c *mockPollsController) findPoll(uuid string) (*models.Poll, error) {
	for _, poll := range c.polls {
		if poll.UUID == uuid {
			return poll, nil
		}
	}
	return nil, errors.New("no poll found")
}

type mockIpsController struct {
	sync.RWMutex
	polls map[string]map[string]struct{}
}

func newMockIpsController() *mockIpsController {
	return &mockIpsController{
		polls: make(map[string]map[string]struct{}),
	}
}

func (c *mockIpsController) AddIPForPoll(uuid string, ip string) error {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.polls[uuid]; !ok {
		c.polls[uuid] = make(map[string]struct{})
	}

	if _, ok := c.polls[uuid][ip]; ok {
		return errors.New("ip already exists")
	}

	c.polls[uuid][ip] = struct{}{}
	return nil
}

func (c *mockIpsController) IsVoteAllowedForIP(uuid string, ip string) bool {
	c.RLock()
	defer c.RUnlock()

	if _, ok := c.polls[uuid]; !ok {
		return true
	}

	_, alreadyVoted := c.polls[uuid][ip]
	return !alreadyVoted
}
