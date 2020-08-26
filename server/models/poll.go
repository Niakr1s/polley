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

// PollArgs is constructor arguments.
type PollArgs struct {
	TimeLimitMinutes int
	AllowMultiple    int
	Name             string
	Choices          []string
}

// normalize sets uninitialized or broken members to appropriate value.
func (args *PollArgs) normalize() {
	if args.TimeLimitMinutes < 0 {
		args.TimeLimitMinutes = 0
	}
	if args.AllowMultiple < 1 {
		args.AllowMultiple = 1
	} else if args.AllowMultiple > len(args.Choices) {
		args.AllowMultiple = len(args.Choices)
	}
	if args.Choices == nil {
		args.Choices = []string{}
	}
}

// NewPoll create new Poll, assigning each choice with zero votes, generating new UUID and setting Expires time.
func NewPoll(args PollArgs) *Poll {
	args.normalize()

	res := NewEmptyPoll()
	res.Expires = time.Now().Add(time.Minute * time.Duration(args.TimeLimitMinutes)).UTC().Round(time.Second)

	for _, text := range args.Choices {
		res.Choices = append(res.Choices, NewChoice(text))
	}

	res.AllowMultiple = args.AllowMultiple

	res.Name = args.Name
	return res
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
