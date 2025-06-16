package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"github.com/yosa/ocr-golang-back/util"
)

func createRandomUser(t *testing.T) User {

	arg := CreateUserParams{
		Username:     util.RandomUsername(),
		Email:        util.RandomEmail(),
		PasswordHash: pgtype.Text{String: util.RandomPasswordHash(), Valid: true},
		Provider:     pgtype.Text{String: util.RandomProvider(), Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)
	require.Equal(t, arg.Provider, user.Provider)

	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {

	arg := CreateUserParams{
		Username:     util.RandomUsername(),
		Email:        util.RandomEmail(),
		PasswordHash: pgtype.Text{String: util.RandomPasswordHash(), Valid: true},
		Provider:     pgtype.Text{String: util.RandomProvider(), Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)
	require.Equal(t, arg.Provider, user.Provider)

	require.NotZero(t, user.CreatedAt)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user1.Username)
	require.NoError(t, err)

	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)

	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, user2)

}

func TestGetUserByEmail(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordHash.String, user2.PasswordHash.String)
	require.Equal(t, user1.Provider.String, user2.Provider.String)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)

}

func TestGetUserByUsername(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordHash.String, user2.PasswordHash.String)
	require.Equal(t, user1.Provider.String, user2.Provider.String)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestListUsers(t *testing.T) {

	for range 10 {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {

		require.NotEmpty(t, user)

	}
}

func TestUpdateUserPassword(t *testing.T) {

	user1 := createRandomUser(t)

	arg := UpdateUserPasswordParams{
		Username:     user1.Username,
		PasswordHash: pgtype.Text{String: util.RandomPasswordHash(), Valid: true},
	}

	err := testQueries.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)

	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.NotEqual(t, user1.PasswordHash.String, user2.PasswordHash.String)
	require.Equal(t, user1.Provider.String, user2.Provider.String)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)

}
