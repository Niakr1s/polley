package pg

import (
	"polley/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPollDB(t *testing.T) {
	pool, err := initTestPool()
	assert.NoError(t, err)
	applyTestMigrations(pool)

	pollDB := NewPollDB(pool)

	poll, err := models.NewPoll(30, []string{"a", "b", "c"})
	assert.NoError(t, err)

	err = pollDB.Create(poll)
	assert.NoError(t, err)

	storedPoll, err := pollDB.Read(poll.UUID)
	assert.NoError(t, err)

	assert.Equal(t, poll, storedPoll)

	const votes = 3
	for i := 0; i < votes; i++ {
		pollDB.Increment(poll.UUID, poll.Choices[0].Text)
	}

	storedPoll, err = pollDB.Read(poll.UUID)
	assert.NoError(t, err)
	assert.Equal(t, votes, storedPoll.Choices[0].Votes)
}