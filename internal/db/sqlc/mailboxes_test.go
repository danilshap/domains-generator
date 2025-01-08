package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/danilshap/domains-generator/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateMailbox(t *testing.T) {
	store := NewStore(testDB)
	user := createRandomUser(t)
	domain := createRandomDomain(t, user)

	arg := CreateMailboxParams{
		Address:  fmt.Sprintf("%s@%s", utils.RandomString(10), domain.Name),
		Password: utils.RandomString(12),
		DomainID: domain.ID,
		UserID:   user.ID,
		Status:   1,
	}

	mailbox, err := store.CreateMailbox(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, mailbox)

	require.Equal(t, arg.Address, mailbox.Address)
	require.Equal(t, arg.Password, mailbox.Password)
	require.Equal(t, arg.DomainID, mailbox.DomainID)
	require.Equal(t, arg.UserID, mailbox.UserID)
	require.Equal(t, arg.Status, mailbox.Status)
	require.NotZero(t, mailbox.ID)
	require.NotZero(t, mailbox.CreatedAt)
}

func TestGetMailboxByID(t *testing.T) {
	store := NewStore(testDB)
	user := createRandomUser(t)
	domain := createRandomDomain(t, user)
	mailbox1 := createRandomMailbox(t, domain, user)

	mailbox2, err := store.GetMailboxByID(context.Background(), mailbox1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, mailbox2)

	require.Equal(t, mailbox1.ID, mailbox2.ID)
	require.Equal(t, mailbox1.Address, mailbox2.Address)
	require.Equal(t, mailbox1.Password, mailbox2.Password)
	require.Equal(t, mailbox1.DomainID, mailbox2.DomainID)
	require.Equal(t, mailbox1.Status, mailbox2.Status)
}

func TestGetMailboxesWithFilters(t *testing.T) {
	store := NewStore(testDB)
	user := createRandomUser(t)
	domain := createRandomDomain(t, user)

	// Create mailboxes with different statuses
	activeMailbox := createRandomMailbox(t, domain, user)

	testCases := []struct {
		name          string
		filters       GetMailboxesWithFiltersParams
		checkResponse func(t *testing.T, mailboxes []GetMailboxesWithFiltersRow)
	}{
		{
			name: "FilterByStatus",
			filters: GetMailboxesWithFiltersParams{
				StatusFilter: []int32{1},
				UserID:       user.ID,
				PageLimit:    10,
			},
			checkResponse: func(t *testing.T, mailboxes []GetMailboxesWithFiltersRow) {
				require.Len(t, mailboxes, 1)
				require.Equal(t, activeMailbox.Address, mailboxes[0].Address)
			},
		},
		{
			name: "FilterByDomain",
			filters: GetMailboxesWithFiltersParams{
				DomainFilter: []int32{domain.ID},
				UserID:       user.ID,
				PageLimit:    10,
			},
			checkResponse: func(t *testing.T, mailboxes []GetMailboxesWithFiltersRow) {
				require.Len(t, mailboxes, 1)
			},
		},
		{
			name: "SearchByAddress",
			filters: GetMailboxesWithFiltersParams{
				SearchQuery: activeMailbox.Address,
				UserID:      user.ID,
				PageLimit:   10,
			},
			checkResponse: func(t *testing.T, mailboxes []GetMailboxesWithFiltersRow) {
				require.Len(t, mailboxes, 1)
				require.Equal(t, activeMailbox.Address, mailboxes[0].Address)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mailboxes, err := store.GetMailboxesWithFilters(context.Background(), tc.filters)
			require.NoError(t, err)
			tc.checkResponse(t, mailboxes)
		})
	}
}

func TestUpdateMailboxStatus(t *testing.T) {
	store := NewStore(testDB)
	user := createRandomUser(t)
	domain := createRandomDomain(t, user)
	mailbox := createRandomMailbox(t, domain, user)

	err := store.UpdateMailboxesStatusByID(context.Background(), UpdateMailboxesStatusByIDParams{
		ID:     mailbox.ID,
		Status: 2,
	})
	require.NoError(t, err)

	updatedMailbox, err := store.GetMailboxByID(context.Background(), mailbox.ID)
	require.NoError(t, err)
	require.Equal(t, int32(2), updatedMailbox.Status)
}

func createRandomMailbox(t *testing.T, domain Domain, user User) Mailbox {
	arg := CreateMailboxParams{
		Address:  utils.RandomEmail(),
		Password: utils.RandomString(12),
		DomainID: domain.ID,
		UserID:   user.ID,
		Status:   1,
	}

	mailbox, err := testStore.CreateMailbox(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, mailbox)

	return mailbox
}
