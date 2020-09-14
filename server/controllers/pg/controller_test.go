package pg

import (
	"polley/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestController(t *testing.T) {
	pool, err := initTestPool()
	assert.NoError(t, err)
	applyTestMigrations(pool)

	pgController := NewPollController(pool)

	// testing PollController interface

	choicesArr := []string{"a", "b", "c"}
	poll := models.NewPoll(models.PollArgs{TimeLimitMinutes: 30, Choices: choicesArr})

	err = pgController.Create(poll)
	assert.NoError(t, err)

	storedPoll, err := pgController.Read(poll.UUID)
	assert.NoError(t, err)

	assert.Equal(t, poll, storedPoll)

	const votes = 3
	for i := 0; i < votes; i++ {
		pgController.Increment(poll.UUID, choicesArr)
	}

	storedPoll, err = pgController.Read(poll.UUID)
	assert.NoError(t, err)
	assert.Equal(t, votes, storedPoll.Choices[0].Votes)
	assert.Equal(t, votes, storedPoll.Choices[1].Votes)
	assert.Equal(t, votes, storedPoll.Choices[2].Votes)

	// testing IPsController interface

	storedUUID := poll.UUID

	err = pgController.AddIPForPoll(storedUUID, "ip")
	assert.NoError(t, err)

	assert.False(t, pgController.IsVoteAllowedForIP(storedUUID, "ip"))
	assert.True(t, pgController.IsVoteAllowedForIP(storedUUID, "other_ip"))

	otherUUID := "some other UUID"
	err = pgController.AddIPForPoll(otherUUID, "ip")
	assert.Error(t, err)

	assert.True(t, pgController.IsVoteAllowedForIP(otherUUID, "ip"))
	assert.True(t, pgController.IsVoteAllowedForIP(otherUUID, "other_ip"))
}
