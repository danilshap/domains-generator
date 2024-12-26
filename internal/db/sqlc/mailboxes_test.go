package db

import (
	"context"
	"testing"

	"github.com/danilshap/domains-generator/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateMailbox(t *testing.T) {
	domain := createRandomDomain(t)
	arg := CreateMailboxParams{
		Address:  utils.RandomEmail(),
		Password: utils.RandomString(12),
		DomainID: domain.ID,
		Status:   1,
	}

	mailbox, err := testQueries.CreateMailbox(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, mailbox)

	require.Equal(t, arg.Address, mailbox.Address)
	require.Equal(t, arg.Password, mailbox.Password)
	require.Equal(t, arg.DomainID, mailbox.DomainID)
	require.Equal(t, arg.Status, mailbox.Status)
	require.False(t, mailbox.IsDeleted.Bool)
	require.NotZero(t, mailbox.ID)
	require.NotZero(t, mailbox.CreatedAt)
}

func TestGetAllMailboxes(t *testing.T) {
	// Create multiple mailboxes
	for i := 0; i < 5; i++ {
		createRandomMailbox(t)
	}

	mailboxes, err := testQueries.GetAllMailboxes(context.Background(), GetAllMailboxesParams{
		Limit:  10,
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, mailboxes)
	require.GreaterOrEqual(t, len(mailboxes), 5)

	for _, mailbox := range mailboxes {
		require.NotEmpty(t, mailbox.ID)
		require.NotEmpty(t, mailbox.Address)
		require.NotZero(t, mailbox.Status)
		require.NotZero(t, mailbox.CreatedAt)
	}
}

func TestGetMailboxesByDomain(t *testing.T) {
	mailbox1 := createRandomMailbox(t)

	mailboxes, err := testQueries.GetMailboxesByDomain(context.Background(), mailbox1.Address)
	require.NoError(t, err)
	require.NotEmpty(t, mailboxes)

	for _, mailbox := range mailboxes {
		require.Equal(t, mailbox1.Address, mailbox.Address)
		require.NotEmpty(t, mailbox.Password)
		require.NotZero(t, mailbox.Status)
		require.NotZero(t, mailbox.CreatedAt)
	}
}

func TestGetMailboxesByDomainName(t *testing.T) {
	mailbox1 := createRandomMailbox(t)

	mailboxes, err := testQueries.GetMailboxesByDomainName(context.Background(), mailbox1.DomainID)
	require.NoError(t, err)
	require.NotEmpty(t, mailboxes)

	found := false
	for _, mailbox := range mailboxes {
		if mailbox.ID == mailbox1.ID {
			found = true
			require.Equal(t, mailbox1.Address, mailbox.Address)
			require.Equal(t, mailbox1.Status, mailbox.Status)
		}
	}
	require.True(t, found)
}

func TestSetMailboxStatus(t *testing.T) {
	mailbox1 := createRandomMailbox(t)

	arg := SetMailboxStatusParams{
		ID:     mailbox1.ID,
		Status: 2,
	}

	err := testQueries.SetMailboxStatus(context.Background(), arg)
	require.NoError(t, err)

	// Check that status was updated
	mailboxes, err := testQueries.GetMailboxesByDomain(context.Background(), mailbox1.Address)
	require.NoError(t, err)
	require.NotEmpty(t, mailboxes)
	require.Equal(t, arg.Status, mailboxes[0].Status)
}

func TestDeleteMailbox(t *testing.T) {
	mailbox1 := createRandomMailbox(t)

	err := testQueries.DeleteMailbox(context.Background(), mailbox1.ID)
	require.NoError(t, err)

	// Check that mailbox was soft deleted
	mailboxes, err := testQueries.GetMailboxesByDomain(context.Background(), mailbox1.Address)
	require.NoError(t, err)
	require.Empty(t, mailboxes)

	// Verify that the mailbox still exists in DB but is marked as deleted
	var deletedMailbox Mailbox
	err = testQueries.db.QueryRowContext(context.Background(), `
		SELECT id, address, password, domain_id, created_at, status, is_deleted 
		FROM mailboxes WHERE id = $1
	`, mailbox1.ID).Scan(
		&deletedMailbox.ID,
		&deletedMailbox.Address,
		&deletedMailbox.Password,
		&deletedMailbox.DomainID,
		&deletedMailbox.CreatedAt,
		&deletedMailbox.Status,
		&deletedMailbox.IsDeleted,
	)
	require.NoError(t, err)
	require.True(t, deletedMailbox.IsDeleted.Bool)
}

func TestGetMailboxesWithFilters(t *testing.T) {
	// Create a test domain
	domain := createRandomDomain(t)

	// Create active and inactive mailboxes
	for i := 0; i < 2; i++ {
		createMailboxWithStatus(t, domain.ID, 1) // active
		createMailboxWithStatus(t, domain.ID, 2) // inactive
	}

	// Test status filter
	activeMailboxes, err := testQueries.GetMailboxesWithFilters(context.Background(), GetMailboxesWithFiltersParams{
		Column1: []int32{1},
		Column2: []int32{},
		Column3: "",
		Limit:   10,
		Offset:  0,
	})
	require.NoError(t, err)
	for _, mailbox := range activeMailboxes {
		require.Equal(t, int32(1), mailbox.Status)
	}

	// Test domain filter
	domainMailboxes, err := testQueries.GetMailboxesWithFilters(context.Background(), GetMailboxesWithFiltersParams{
		Column1: []int32{},
		Column2: []int32{domain.ID},
		Column3: "",
		Limit:   10,
		Offset:  0,
	})
	require.NoError(t, err)
	for _, mailbox := range domainMailboxes {
		require.Equal(t, domain.ID, mailbox.DomainID)
	}
}

func TestGetMailboxesStats(t *testing.T) {
	// Create a test domain
	domain := createRandomDomain(t)

	// Create mailboxes with different statuses
	createMailboxWithStatus(t, domain.ID, 1)
	createMailboxWithStatus(t, domain.ID, 1)
	createMailboxWithStatus(t, domain.ID, 2)

	// Get statistics
	stats, err := testQueries.GetMailboxesStats(context.Background())
	require.NoError(t, err)

	// Check statistics
	require.GreaterOrEqual(t, stats.TotalCount, int64(3))
	require.GreaterOrEqual(t, stats.ActiveCount, int64(2))
	require.GreaterOrEqual(t, stats.InactiveCount, int64(1))
	require.GreaterOrEqual(t, stats.DomainsCount, int64(1))
}

// Helper function to create a mailbox with specific status
func createMailboxWithStatus(t *testing.T, domainID int32, status int32) Mailbox {
	arg := CreateMailboxParams{
		Address:  utils.RandomEmail(),
		Password: utils.RandomString(12),
		DomainID: domainID,
		Status:   status,
	}

	mailbox, err := testQueries.CreateMailbox(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, mailbox)

	return mailbox
}

// Helper function to create a random mailbox
func createRandomMailbox(t *testing.T) Mailbox {
	domain := createRandomDomain(t)

	arg := CreateMailboxParams{
		Address:  utils.RandomEmail(),
		Password: utils.RandomString(12),
		DomainID: domain.ID,
		Status:   1,
	}

	mailbox, err := testQueries.CreateMailbox(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, mailbox)

	return mailbox
}
