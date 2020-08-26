package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Poll represents poll.
type Poll struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`

	Expires       time.Time `json:"expires"`
	AllowMultiple int       `json:"allowMultiple"`

	// Choices - Votes pairs.
	Choices []Choice `json:"choices"`
}

// Choice represents single choice.
type Choice struct {
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

// Errors
var (
	ErrWrongTimeLimit = errors.New("Wrong timeLimit argument")
)

// NewPoll create new Poll, assigning each choice with zero votes, generating new UUID and setting Expires time.
func NewPoll(timeLimitMinutes int, choices []string) (*Poll, error) {
	if timeLimitMinutes < 0 {
		return nil, ErrWrongTimeLimit
	}
	res := NewEmptyPoll()
	res.Expires = time.Now().Add(time.Minute * time.Duration(timeLimitMinutes)).UTC().Round(time.Second)

	for _, text := range choices {
		res.Choices = append(res.Choices, NewChoice(text))
	}
	return res, nil
}

// WithAllowMultiple sets allowMultiple, bounding it to [1..len(p.Choices)]
func (p *Poll) WithAllowMultiple(allowMultiple int) *Poll {
	if allowMultiple < 1 {
		allowMultiple = 1
	} else if allowMultiple > len(p.Choices) {
		allowMultiple = len(p.Choices)
	}
	p.AllowMultiple = allowMultiple
	return p
}

// WithName sets name.
func (p *Poll) WithName(name string) *Poll {
	p.Name = name
	return p
}

// NewEmptyPoll creates new Poll, with default values.
func NewEmptyPoll() *Poll {
	return &Poll{
		UUID:          uuid.New().String(),
		Choices:       make([]Choice, 0),
		AllowMultiple: 1,
	}
}

// IsExpired returns expired status of poll.
func (p *Poll) IsExpired() bool {
	return p.WillBeExpiredAt(time.Now())
}

// WillBeExpiredAt returns if poll will be expired at point of time.
func (p *Poll) WillBeExpiredAt(t time.Time) bool {
	return p.Expires.Before(t)
}

func (p *Poll) String() string {
	return fmt.Sprintf("Poll: { UUID: %s, expires: %v, AllowMultiple: %v, choices: %v }", p.UUID, p.Expires, p.AllowMultiple, p.Choices)
}

// NewChoice constructs new Choice with zero votes.
func NewChoice(text string) Choice {
	return Choice{Text: text, Votes: 0}
}
