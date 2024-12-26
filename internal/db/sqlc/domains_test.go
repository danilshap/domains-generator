package db

import (
	"context"
	"testing"

	"github.com/danilshap/domains-generator/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateDomain(t *testing.T) {
	arg := CreateDomainParams{
		Name:     utils.RandomName(),
		Provider: utils.RandomProvider(),
		Status:   1,
	}

	domain, err := testQueries.CreateDomain(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, domain)

	require.Equal(t, arg.Name, domain.Name)
	require.Equal(t, arg.Provider, domain.Provider)
	require.Equal(t, arg.Status, domain.Status)
	require.False(t, domain.IsDeleted.Bool)
	require.NotZero(t, domain.ID)
	require.NotZero(t, domain.CreatedAt)
}

func TestGetDomainByName(t *testing.T) {
	// First create a domain
	createdDomain := createRandomDomain(t)

	// Get domain by name
	domain, err := testQueries.GetDomainByName(context.Background(), createdDomain.Name)
	require.NoError(t, err)
	require.NotEmpty(t, domain)

	require.Equal(t, createdDomain.ID, domain.ID)
	require.Equal(t, createdDomain.Name, domain.Name)
	require.Equal(t, createdDomain.Provider, domain.Provider)
	require.Equal(t, createdDomain.Status, domain.Status)
	require.Equal(t, createdDomain.IsDeleted, domain.IsDeleted)
}

func TestGetAllDomains(t *testing.T) {
	// Create multiple domains
	for i := 0; i < 5; i++ {
		createRandomDomain(t)
	}

	domains, err := testQueries.GetAllDomains(context.Background(), GetAllDomainsParams{
		Limit:  5,
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, domains)
	require.GreaterOrEqual(t, len(domains), 5)

	for _, domain := range domains {
		require.NotEmpty(t, domain.ID)
		require.NotEmpty(t, domain.Name)
		require.NotEmpty(t, domain.Provider)
		require.NotZero(t, domain.Status)
		require.NotZero(t, domain.CreatedAt)
	}
}

func TestSetDomainStatus(t *testing.T) {
	domain := createRandomDomain(t)

	arg := SetDomainStatusParams{
		ID:     domain.ID,
		Status: 2,
	}

	err := testQueries.SetDomainStatus(context.Background(), arg)
	require.NoError(t, err)

	// Check that status was updated
	updatedDomain, err := testQueries.GetDomainByName(context.Background(), domain.Name)
	require.NoError(t, err)
	require.Equal(t, arg.Status, updatedDomain.Status)
}

func TestDeleteDomain(t *testing.T) {
	domain := createRandomDomain(t)

	err := testQueries.DeleteDomain(context.Background(), domain.ID)
	require.NoError(t, err)

	// Check that domain was soft deleted
	_, err = testQueries.GetDomainByName(context.Background(), domain.Name)
	require.Error(t, err)

	// Verify that the domain still exists in DB but is marked as deleted
	var deletedDomain Domain
	err = testQueries.db.QueryRowContext(context.Background(), `
		SELECT id, name, provider, status, created_at, expires_at, is_deleted, settings 
		FROM domains WHERE id = $1
	`, domain.ID).Scan(
		&deletedDomain.ID,
		&deletedDomain.Name,
		&deletedDomain.Provider,
		&deletedDomain.Status,
		&deletedDomain.CreatedAt,
		&deletedDomain.ExpiresAt,
		&deletedDomain.IsDeleted,
		&deletedDomain.Settings,
	)
	require.NoError(t, err)
	require.True(t, deletedDomain.IsDeleted.Bool)
}

func TestUpdateDomainAndMailboxesStatus(t *testing.T) {
	domain := createRandomDomain(t)

	for i := 0; i < 3; i++ {
		arg := CreateMailboxParams{
			Address:  utils.RandomEmail(),
			Password: utils.RandomString(12),
			DomainID: domain.ID,
			Status:   1,
		}
		_, err := testQueries.CreateMailbox(context.Background(), arg)
		require.NoError(t, err)
	}

	err := testStore.UpdateDomainAndMailboxesStatus(context.Background(), domain.ID, 2)
	require.NoError(t, err)

	updatedDomain, err := testQueries.GetDomainByID(context.Background(), domain.ID)
	require.NoError(t, err)
	require.Equal(t, int32(2), updatedDomain.Status)

	mailboxes, err := testQueries.GetMailboxesByDomainID(context.Background(), GetMailboxesByDomainIDParams{
		DomainID: domain.ID,
		Limit:    10,
		Offset:   0,
	})
	require.NoError(t, err)
	for _, mailbox := range mailboxes {
		require.Equal(t, int32(2), mailbox.Status)
	}
}

func TestGetDomainsWithPagination(t *testing.T) {
	// Create multiple domains
	for i := 0; i < 5; i++ {
		createRandomDomain(t)
	}

	// Test first page
	page1, err := testQueries.GetAllDomains(context.Background(), GetAllDomainsParams{
		Limit:  3,
		Offset: 0,
	})
	require.NoError(t, err)
	require.Len(t, page1, 3)

	// Test second page
	page2, err := testQueries.GetAllDomains(context.Background(), GetAllDomainsParams{
		Limit:  3,
		Offset: 3,
	})
	require.NoError(t, err)
	require.NotEmpty(t, page2)

	// Check that IDs are not repeated between pages
	page1IDs := make(map[int32]bool)
	for _, domain := range page1 {
		page1IDs[domain.ID] = true
	}
	for _, domain := range page2 {
		require.False(t, page1IDs[domain.ID])
	}
}

// Helper function to create a random domain
func createRandomDomain(t *testing.T) Domain {
	arg := CreateDomainParams{
		Name:     utils.RandomName(),
		Provider: utils.RandomProvider(),
		Status:   1,
	}

	domain, err := testQueries.CreateDomain(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, domain)

	return domain
}
