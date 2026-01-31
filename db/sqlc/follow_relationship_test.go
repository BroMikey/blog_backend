package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func followInList(items []FollowRelationship, followerID, followedID int64) bool {
	for _, item := range items {
		if item.FollowerUid == followerID && item.FollowedUid == followedID {
			return true
		}
	}
	return false
}

func TestCreateFollowRelationship(t *testing.T) {
	follower := createRandomUser(t)
	followed := createRandomUser(t)
	defer cleanupUser(t, follower.Uid)
	defer cleanupUser(t, followed.Uid)

	created, err := testQueries.CreateFollowRelationship(context.Background(), CreateFollowRelationshipParams{FollowerUid: follower.Uid, FollowedUid: followed.Uid})
	require.NoError(t, err)
	require.Equal(t, follower.Uid, created.FollowerUid)
	require.Equal(t, followed.Uid, created.FollowedUid)

	err = testQueries.DeleteFollowRelationship(context.Background(), DeleteFollowRelationshipParams{FollowerUid: follower.Uid, FollowedUid: followed.Uid})
	require.NoError(t, err)
}

func TestGetFollowRelationship(t *testing.T) {
	follower := createRandomUser(t)
	followed := createRandomUser(t)
	defer cleanupUser(t, follower.Uid)
	defer cleanupUser(t, followed.Uid)

	created, err := testQueries.CreateFollowRelationship(context.Background(), CreateFollowRelationshipParams{FollowerUid: follower.Uid, FollowedUid: followed.Uid})
	require.NoError(t, err)

	fetched, err := testQueries.GetFollowRelationship(context.Background(), GetFollowRelationshipParams{FollowerUid: follower.Uid, FollowedUid: followed.Uid})
	require.NoError(t, err)
	require.Equal(t, created.ID, fetched.ID)
	require.Equal(t, follower.Uid, fetched.FollowerUid)
	require.Equal(t, followed.Uid, fetched.FollowedUid)
}

func TestGetFollowRelationshipByID(t *testing.T) {
	follower := createRandomUser(t)
	followed := createRandomUser(t)
	defer cleanupUser(t, follower.Uid)
	defer cleanupUser(t, followed.Uid)

	created, err := testQueries.CreateFollowRelationship(context.Background(), CreateFollowRelationshipParams{FollowerUid: follower.Uid, FollowedUid: followed.Uid})
	require.NoError(t, err)

	fetched, err := testQueries.GetFollowRelationshipByID(context.Background(), created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, fetched.ID)
}

func TestListFollowersAndFollowing(t *testing.T) {
	follower := createRandomUser(t)
	followed := createRandomUser(t)
	defer cleanupUser(t, follower.Uid)
	defer cleanupUser(t, followed.Uid)

	_, err := testQueries.CreateFollowRelationship(context.Background(), CreateFollowRelationshipParams{FollowerUid: follower.Uid, FollowedUid: followed.Uid})
	require.NoError(t, err)

	following, err := testQueries.ListFollowing(context.Background(), ListFollowingParams{FollowerUid: follower.Uid, Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.True(t, followInList(following, follower.Uid, followed.Uid))

	followers, err := testQueries.ListFollowers(context.Background(), ListFollowersParams{FollowedUid: followed.Uid, Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.True(t, followInList(followers, follower.Uid, followed.Uid))
}

func TestDeleteFollowRelationship(t *testing.T) {
	follower := createRandomUser(t)
	followed := createRandomUser(t)
	defer cleanupUser(t, follower.Uid)
	defer cleanupUser(t, followed.Uid)

	created, err := testQueries.CreateFollowRelationship(context.Background(), CreateFollowRelationshipParams{FollowerUid: follower.Uid, FollowedUid: followed.Uid})
	require.NoError(t, err)

	err = testQueries.DeleteFollowRelationship(context.Background(), DeleteFollowRelationshipParams{FollowerUid: follower.Uid, FollowedUid: followed.Uid})
	require.NoError(t, err)

	_, err = testQueries.GetFollowRelationshipByID(context.Background(), created.ID)
	require.Error(t, err)
}
