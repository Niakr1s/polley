package storage

import (
	"polley/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mockPollsController(t *testing.T) {
	c := newMockPollsController()

	poll := models.NewEmptyPoll()
	poll.Choices = append(poll.Choices, models.NewChoice("a"), models.NewChoice("b"), models.NewChoice("c"))

	err := c.Create(poll)
	assert.NoError(t, err)

	err = c.Create(poll)
	assert.Error(t, err)

	err = c.Create(models.NewEmptyPoll())
	assert.NoError(t, err)

	storedPoll, err := c.Read(poll.UUID)
	assert.NoError(t, err)
	assert.Equal(t, poll.UUID, storedPoll.UUID)

	_, err = c.Read("someOtherUUID")
	assert.Error(t, err)

	err = c.Increment(poll.UUID, []string{"a", "b"})
	assert.NoError(t, err)
	storedPoll, err = c.Read(poll.UUID)
	assert.NoError(t, err)
	assert.Equal(t, 1, storedPoll.Choices[0].Votes)
	assert.Equal(t, 1, storedPoll.Choices[1].Votes)
	assert.Equal(t, 0, storedPoll.Choices[2].Votes)

	err = c.Increment(poll.UUID, []string{"a", "b", "f"})
	assert.Error(t, err)
	storedPoll, err = c.Read(poll.UUID)
	assert.NoError(t, err)
	// asserting not incremented, because of error
	assert.Equal(t, 1, storedPoll.Choices[0].Votes)
	assert.Equal(t, 1, storedPoll.Choices[1].Votes)
	assert.Equal(t, 0, storedPoll.Choices[2].Votes)

	uuids, err := c.GetNPollsUUIDs(10, 0)
	assert.NoError(t, err)
	assert.Len(t, uuids, 2)

	_, err = c.GetNPollsUUIDs(5, 2)
	assert.Error(t, err)

	total := c.GetTotal()
	assert.Equal(t, 2, total)
}

func Test_mockIpsController(t *testing.T) {
	c := newMockIpsController()

	uuid1 := "uuid1"
	uuid2 := "uuid2"
	ip1 := "111.111.111.111"
	ip2 := "222.222.222.222"
	var err error

	err = c.AddIPForPoll(uuid1, ip1)
	assert.NoError(t, err)
	err = c.AddIPForPoll(uuid1, ip1)
	assert.Error(t, err)

	isAllowed1 := c.IsVoteAllowedForIP(uuid1, ip1)
	isAllowed2 := c.IsVoteAllowedForIP(uuid1, ip2)
	assert.False(t, isAllowed1)
	assert.True(t, isAllowed2)

	isAllowed1 = c.IsVoteAllowedForIP(uuid2, ip1)
	isAllowed2 = c.IsVoteAllowedForIP(uuid2, ip2)
	assert.True(t, isAllowed1)
	assert.True(t, isAllowed2)

}
