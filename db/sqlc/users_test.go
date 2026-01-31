package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/BroMikey/blog_backend/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) Users {
	t.Helper()

	args := CreateUserParams{
		Username:     fmt.Sprintf("user_%s_%d", utils.RandomString(6), utils.RandomInt(0, 1_000_000_000)),
		Email:        utils.RandomEmail(),
		PasswordHash: utils.RandomPassword(),
	}

	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.Email, user.Email)
	require.Equal(t, args.PasswordHash, user.PasswordHash)

	require.NotZero(t, user.Uid)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)
	require.Equal(t, int16(1), user.Status)

	return user
}

// cleanupUser deletes a user used in tests
func cleanupUser(t *testing.T, uid int64) {
	t.Helper()

	err := testQueries.DeleteUser(context.Background(), uid)
	require.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	user := createRandomUser(t)
	defer cleanupUser(t, user.Uid)

	// Verify the user was created with correct initial values
	require.NotZero(t, user.Uid)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)
	require.Equal(t, int16(1), user.Status)
}

func TestGetUser(t *testing.T) {
	createdUser := createRandomUser(t)
	defer cleanupUser(t, createdUser.Uid)

	fetchedUser, err := testQueries.GetUser(context.Background(), createdUser.Uid)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedUser)

	require.Equal(t, createdUser.Uid, fetchedUser.Uid)
	require.Equal(t, createdUser.Username, fetchedUser.Username)
	require.Equal(t, createdUser.Email, fetchedUser.Email)
	require.Equal(t, createdUser.PasswordHash, fetchedUser.PasswordHash)
}

func TestListUsers(t *testing.T) {
	args := ListUsersParams{
		Limit:  5,
		Offset: 0,
	}

	// Create multiple test users for pagination testing
	for i := 0; i < 10; i++ {
		user := createRandomUser(t)
		defer cleanupUser(t, user.Uid)
	}

	// Call ListUsers for pagination test
	users, err := testQueries.ListUsers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, users, int(args.Limit))

	for _, user := range users {
		require.NotEmpty(t, user)
		require.NotZero(t, user.Uid)
		require.NotEmpty(t, user.Username)
		require.NotEmpty(t, user.Email)
		require.NotEmpty(t, user.PasswordHash)
		require.NotZero(t, user.CreatedAt)
		require.NotZero(t, user.UpdatedAt)
		require.Equal(t, int16(1), user.Status)
	}
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)
	defer cleanupUser(t, user.Uid)

	// Update email and bio
	newEmail := utils.RandomEmail()
	newBio := "Updated bio"
	updateArgs := UpdateUserParams{
		Uid:          user.Uid,
		Username:     user.Username,
		Email:        newEmail,
		PasswordHash: user.PasswordHash,
		Bio:          sql.NullString{String: newBio, Valid: true},
		Avatar:       sql.NullString{Valid: false},
		Status:       user.Status,
	}

	updated, err := testQueries.UpdateUser(context.Background(), updateArgs)
	require.NoError(t, err)
	require.NotEmpty(t, updated)

	// Verify updated fields
	require.Equal(t, newEmail, updated.Email)
	require.Equal(t, newBio, updated.Bio.String)

	// Verify unchanged fields
	require.Equal(t, user.Username, updated.Username)
	require.Equal(t, user.PasswordHash, updated.PasswordHash)
	require.Equal(t, user.Status, updated.Status)
}
