package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/danilshap/domains-generator/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateDomain(t *testing.T) {
	store := NewStore(testDB)
	user := createRandomUser(t)

	arg := CreateDomainParams{
		Name:     utils.RandomString(10),
		Provider: "test",
		Status:   1,
		UserID:   user.ID,
		ExpiresAt: sql.NullTime{
			Time:  time.Now().Add(24 * time.Hour),
			Valid: true,
		},
	}

	domain, err := store.CreateDomain(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, domain)

	require.Equal(t, arg.Name, domain.Name)
	require.Equal(t, arg.Provider, domain.Provider)
	require.Equal(t, arg.Status, domain.Status)
	require.Equal(t, arg.UserID, domain.UserID)
	require.NotZero(t, domain.ID)
	require.NotZero(t, domain.CreatedAt)
}

func TestGetDomainByID(t *testing.T) {
	store := NewStore(testDB)
	user := createRandomUser(t)

	domain1 := createRandomDomain(t, user)
	domain2, err := store.GetDomainByID(context.Background(), domain1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, domain2)

	require.Equal(t, domain1.ID, domain2.ID)
	require.Equal(t, domain1.Name, domain2.Name)
	require.Equal(t, domain1.Provider, domain2.Provider)
	require.Equal(t, domain1.Status, domain2.Status)
}

func TestGetAllDomains(t *testing.T) {
	store := NewStore(testDB)
	user := createRandomUser(t)

	// Create multiple domains
	for i := 0; i < 5; i++ {
		createRandomDomain(t, user)
	}

	arg := GetAllDomainsParams{
		Limit:  5,
		Offset: 0,
		UserID: user.ID,
	}

	domains, err := store.GetAllDomains(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, domains)
	require.Len(t, domains, 5)

	for _, domain := range domains {
		require.NotEmpty(t, domain)
		require.Equal(t, user.ID, domain.UserID)
	}
}

func createRandomDomain(t *testing.T, user User) Domain {
	arg := CreateDomainParams{
		Name:     utils.RandomString(10),
		Provider: "test",
		Status:   1,
		UserID:   user.ID,
		ExpiresAt: sql.NullTime{
			Time:  time.Now().Add(24 * time.Hour),
			Valid: true,
		},
	}

	domain, err := testStore.CreateDomain(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, domain)

	return domain
}
