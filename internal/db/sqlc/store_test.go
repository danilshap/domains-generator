package db

import (
	"context"
	"testing"
	"time"

	"github.com/danilshap/domains-generator/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestUpdateDomainAndMailboxesStatusTx(t *testing.T) {
	store := testStore

	domain := createRandomDomain(t)

	var mailboxes []Mailbox
	for i := 0; i < 5; i++ {
		arg := CreateMailboxParams{
			Address:  utils.RandomEmail(),
			Password: utils.RandomString(12),
			DomainID: domain.ID,
			Status:   1,
		}
		mailbox, err := store.CreateMailbox(context.Background(), arg)
		require.NoError(t, err)
		mailboxes = append(mailboxes, mailbox)
	}

	// Test successful update
	err := store.UpdateDomainAndMailboxesStatus(context.Background(), domain.ID, 2)
	require.NoError(t, err)

	// Check domain status
	updatedDomain, err := store.GetDomainByID(context.Background(), domain.ID)
	require.NoError(t, err)
	require.Equal(t, int32(2), updatedDomain.Status)

	// Check status of all mailboxes
	for _, mailbox := range mailboxes {
		updatedMailbox, err := store.GetMailboxByID(context.Background(), mailbox.ID)
		require.NoError(t, err)
		require.Equal(t, int32(2), updatedMailbox.Status)
	}

	// Test activation
	err = store.UpdateDomainAndMailboxesStatus(context.Background(), domain.ID, 1)
	require.NoError(t, err)

	// Check that domain status changed to active
	updatedDomain, err = store.GetDomainByID(context.Background(), domain.ID)
	require.NoError(t, err)
	require.Equal(t, int32(1), updatedDomain.Status)

	// Mailbox statuses should not change when domain is activated
	for _, mailbox := range mailboxes {
		updatedMailbox, err := store.GetMailboxByID(context.Background(), mailbox.ID)
		require.NoError(t, err)
		require.Equal(t, int32(2), updatedMailbox.Status)
	}
}

func TestConcurrentDomainStatusUpdate(t *testing.T) {
	store := testStore

	// Create test domain
	domain := createRandomDomain(t)

	// Create multiple mailboxes
	for i := 0; i < 5; i++ {
		arg := CreateMailboxParams{
			Address:  utils.RandomEmail(),
			Password: utils.RandomString(12),
			DomainID: domain.ID,
			Status:   1,
		}
		_, err := store.CreateMailbox(context.Background(), arg)
		require.NoError(t, err)
	}

	// Number of concurrent updates
	n := 5
	errs := make(chan error, n)

	// Start concurrent updates
	for i := 0; i < n; i++ {
		status := int32(i%2 + 1) // Alternate between statuses 1 and 2
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			errs <- store.UpdateDomainAndMailboxesStatus(ctx, domain.ID, status)
		}()
	}

	// Check that all updates completed successfully
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// Check final state
	updatedDomain, err := store.GetDomainByID(context.Background(), domain.ID)
	require.NoError(t, err)
	require.Contains(t, []int32{1, 2}, updatedDomain.Status)

	// If domain is inactive, check that all mailboxes are also inactive
	if updatedDomain.Status == 2 {
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
}
