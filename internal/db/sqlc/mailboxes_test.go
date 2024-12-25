package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/danilshap/domains-generator/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateMailbox(t *testing.T) {
	domain := createRandomDomain(t)
	arg := CreateMailboxParams{
		Address:  utils.RandomEmail(),
		Password: utils.RandomString(12),
		DomainID: sql.NullInt32{Int32: domain.ID, Valid: true},
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
	// Создаем несколько почтовых ящиков
	for i := 0; i < 5; i++ {
		createRandomMailbox(t)
	}

	mailboxes, err := testQueries.GetAllMailboxes(context.Background())
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

	// Проверяем что статус обновился через получение почтового ящика
	mailboxes, err := testQueries.GetMailboxesByDomain(context.Background(), mailbox1.Address)
	require.NoError(t, err)
	require.NotEmpty(t, mailboxes)
	require.Equal(t, arg.Status, mailboxes[0].Status)
}

func TestDeleteMailbox(t *testing.T) {
	mailbox1 := createRandomMailbox(t)

	err := testQueries.DeleteMailbox(context.Background(), mailbox1.ID)
	require.NoError(t, err)

	// Проверяем что почтовый ящик удален
	mailboxes, err := testQueries.GetMailboxesByDomain(context.Background(), mailbox1.Address)
	require.NoError(t, err)
	require.Empty(t, mailboxes)
}

// Вспомогательная функция для создания случайного почтового ящика
func createRandomMailbox(t *testing.T) Mailbox {
	domain := createRandomDomain(t)

	arg := CreateMailboxParams{
		Address:  utils.RandomEmail(),
		Password: utils.RandomString(12),
		DomainID: sql.NullInt32{Int32: domain.ID, Valid: true},
		Status:   1,
	}

	mailbox, err := testQueries.CreateMailbox(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, mailbox)

	return mailbox
}
