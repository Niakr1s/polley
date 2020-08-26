package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPoll(t *testing.T) {
	poll := NewPoll(PollArgs{TimeLimitMinutes: 10, Choices: []string{"1", "2", "3"}})

	if poll.UUID == "" {
		t.Errorf("uuid is empty")
	}
	if l := len(poll.Choices); l != 3 {
		t.Errorf("len of poll.Choices = %v, should = %v", l, 3)
	}
	if expires := time.Now().Add(time.Second * 9); expires.After(poll.Expires) {
		t.Errorf("poll.Expires error")
	}
}

func TestPoll_IsExpired(t *testing.T) {
	tests := []struct {
		name            string
		durationMinutes int
		want            bool
	}{
		{"expired", 0, true},
		{"not expired", 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := NewPoll(PollArgs{
				TimeLimitMinutes: tt.durationMinutes,
				Choices:          []string{},
			})
			<-time.After(time.Millisecond * 600) // to bypass rounding
			assert.Equal(t, tt.want, p.IsExpired())
		})
	}
}

func TestPoll_WillBeExpiredAt(t *testing.T) {
	tests := []struct {
		name            string
		durationMinutes int
		at              time.Time
		want            bool
	}{
		{"now second after", 0, time.Now().Add(time.Second), true},
		{"3 minutes at 4 minutes", 3, time.Now().Add(time.Minute * 4), true},
		{"3 minutes at 2 minutes", 3, time.Now().Add(time.Minute * 2), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPoll(PollArgs{TimeLimitMinutes: tt.durationMinutes, Choices: []string{}})
			assert.Equal(t, tt.want, p.WillBeExpiredAt(tt.at))
		})
	}
}
