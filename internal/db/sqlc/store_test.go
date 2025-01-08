package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateDomainAndMailboxesStatus(t *testing.T) {
	store := NewStore(testDB)
	user := createRandomUser(t)

	// Create test domain
	domain := createRandomDomain(t, user)

	// Create test mailboxes
	for i := 0; i < 3; i++ {
		createRandomMailbox(t, domain, user)
	}

	// Test updating status
	err := store.UpdateDomainAndMailboxesStatus(context.Background(), domain.ID, 2)
	require.NoError(t, err)

	// Verify domain status
	updatedDomain, err := store.GetDomainByID(context.Background(), domain.ID)
	require.NoError(t, err)
	require.Equal(t, int32(2), updatedDomain.Status)

	// Verify mailboxes status
	mailboxes, err := store.GetMailboxesByDomainID(context.Background(), GetMailboxesByDomainIDParams{
		DomainID: domain.ID,
		Limit:    10,
		Offset:   0,
	})
	require.NoError(t, err)
	for _, mailbox := range mailboxes {
		require.Equal(t, int32(2), mailbox.Status)
	}
}
