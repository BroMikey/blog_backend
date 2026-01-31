package db

import (
	"context"
	"testing"

	"github.com/BroMikey/blog_backend/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTag(t *testing.T) Tag {
	t.Helper()

	name := utils.RandomString(10)
	tag, err := testQueries.CreateTag(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, tag)
	require.Equal(t, name, tag.Name)
	require.NotZero(t, tag.ID)
	require.NotZero(t, tag.CreatedAt)

	return tag
}

func cleanupTag(t *testing.T, id int64) {
	t.Helper()

	err := testQueries.DeleteTag(context.Background(), id)
	require.NoError(t, err)
}

func tagInList(tags []Tag, id int64) bool {
	for _, tag := range tags {
		if tag.ID == id {
			return true
		}
	}
	return false
}

func TestCreateTag(t *testing.T) {
	tag := createRandomTag(t)
	defer cleanupTag(t, tag.ID)

	require.NotZero(t, tag.ID)
}

func TestGetTagByID(t *testing.T) {
	tag := createRandomTag(t)
	defer cleanupTag(t, tag.ID)

	fetched, err := testQueries.GetTagByID(context.Background(), tag.ID)
	require.NoError(t, err)
	require.Equal(t, tag.ID, fetched.ID)
	require.Equal(t, tag.Name, fetched.Name)
}

func TestGetTagByName(t *testing.T) {
	tag := createRandomTag(t)
	defer cleanupTag(t, tag.ID)

	fetched, err := testQueries.GetTagByName(context.Background(), tag.Name)
	require.NoError(t, err)
	require.Equal(t, tag.ID, fetched.ID)
	require.Equal(t, tag.Name, fetched.Name)
}

func TestListTags(t *testing.T) {
	tag1 := createRandomTag(t)
	defer cleanupTag(t, tag1.ID)
	tag2 := createRandomTag(t)
	defer cleanupTag(t, tag2.ID)

	tags, err := testQueries.ListTags(context.Background(), ListTagsParams{Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.NotEmpty(t, tags)
	require.True(t, tagInList(tags, tag1.ID))
	require.True(t, tagInList(tags, tag2.ID))
}

func TestUpdateTag(t *testing.T) {
	tag := createRandomTag(t)
	defer cleanupTag(t, tag.ID)

	newName := utils.RandomString(12)
	updated, err := testQueries.UpdateTag(context.Background(), UpdateTagParams{ID: tag.ID, Name: newName})
	require.NoError(t, err)
	require.Equal(t, tag.ID, updated.ID)
	require.Equal(t, newName, updated.Name)
}

func TestDeleteTag(t *testing.T) {
	tag := createRandomTag(t)
	err := testQueries.DeleteTag(context.Background(), tag.ID)
	require.NoError(t, err)

	_, err = testQueries.GetTagByID(context.Background(), tag.ID)
	require.Error(t, err)
}
