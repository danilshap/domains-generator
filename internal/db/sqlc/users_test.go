package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/danilshap/domains-generator/pkg/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       utils.RandomString(6),
		Email:          utils.RandomEmail(),
		HashedPassword: utils.RandomString(12),
		FullName:       sql.NullString{String: utils.RandomString(10), Valid: true},
		Role:           "user",
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Role, user.Role)
	require.True(t, user.IsActive)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testStore.GetUserByID(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)

	arg := UpdateUserParams{
		ID:       user1.ID,
		Username: sql.NullString{String: utils.RandomString(6), Valid: true},
		Email:    sql.NullString{String: utils.RandomEmail(), Valid: true},
		FullName: sql.NullString{String: utils.RandomString(10), Valid: true},
		Role:     sql.NullString{String: "admin", Valid: true},
	}

	user2, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, arg.Username.String, user2.Username)
	require.Equal(t, arg.Email.String, user2.Email)
	require.Equal(t, arg.FullName.String, user2.FullName.String)
	require.Equal(t, arg.Role.String, user2.Role)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)

	err := testStore.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testStore.GetUserByID(context.Background(), user1.ID)
	require.Error(t, err)
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testStore.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func TestGetUserByEmail(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testStore.GetUserByEmail(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Email, user2.Email)
}

func TestGetUserByUsername(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testStore.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
}

func TestDeactivateUser(t *testing.T) {
	user1 := createRandomUser(t)

	err := testStore.DeactivateUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testStore.GetUserByID(context.Background(), user1.ID)
	require.Error(t, err) // Should error because GetUserByID only returns active users
	require.Empty(t, user2)
}

func TestUpdateUserPassword(t *testing.T) {
	user1 := createRandomUser(t)
	newPassword := utils.RandomString(12)

	arg := UpdateUserPasswordParams{
		ID:             user1.ID,
		HashedPassword: newPassword,
	}

	err := testStore.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)

	user2, err := testStore.GetUserByID(context.Background(), user1.ID)
	require.NoError(t, err)
	require.Equal(t, newPassword, user2.HashedPassword)
}

func TestGetUserStats(t *testing.T) {
	user := createRandomUser(t)

	stats, err := testStore.GetUserStats(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotNil(t, stats)
	// Initially should have no domains or mailboxes
	require.Equal(t, int64(0), stats.DomainsCount)
	require.Equal(t, int64(0), stats.MailboxesCount)
}
