package pg

import (
	"polley/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	pool, err := initTestPool()
	assert.NoError(t, err)
	applyTestMigrations(pool)

	pgDB := NewDB(pool)

	// testing PollDB interface

	poll := models.NewPoll(models.PollArgs{TimeLimitMinutes: 30, Choices: []string{"a", "b", "c"}})

	err = pgDB.Create(poll)
	assert.NoError(t, err)

	storedPoll, err := pgDB.Read(poll.UUID)
	assert.NoError(t, err)

	assert.Equal(t, poll, storedPoll)

	const votes = 3
	for i := 0; i < votes; i++ {
		pgDB.Increment(poll.UUID, poll.Choices[0].Text)
	}

	storedPoll, err = pgDB.Read(poll.UUID)
	assert.NoError(t, err)
	assert.Equal(t, votes, storedPoll.Choices[0].Votes)

	// testing IPsDB interface

	storedUUID := poll.UUID

	err = pgDB.AddIPForPoll(storedUUID, "ip")
	assert.NoError(t, err)

	assert.False(t, pgDB.IsVoteAllowedForIP(storedUUID, "ip"))
	assert.True(t, pgDB.IsVoteAllowedForIP(storedUUID, "other_ip"))

	otherUUID := "some other UUID"
	err = pgDB.AddIPForPoll(otherUUID, "ip")
	assert.Error(t, err)

	assert.True(t, pgDB.IsVoteAllowedForIP(otherUUID, "ip"))
	assert.True(t, pgDB.IsVoteAllowedForIP(otherUUID, "other_ip"))
}
