package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// article id, tag id, article_tag id
func TestCreateArticleTag(t *testing.T) {
	article := createRandomArticle(t)
	tag := createRandomTag(t)
	defer cleanupArticle(t, article.ID)
	defer cleanupTag(t, tag.ID)

	created, err := testQueries.CreateArticleTag(context.Background(), CreateArticleTagParams{
		ArticleID: article.ID,
		TagID:     tag.ID,
	})
	require.NoError(t, err)
	require.Equal(t, article.ID, created.ArticleID)
	require.Equal(t, tag.ID, created.TagID)

	err = testQueries.DeleteArticleTag(context.Background(), DeleteArticleTagParams{
		ArticleID: article.ID,
		TagID:     tag.ID,
	})
	require.NoError(t, err)
}

// TestGetArticleTagByID tests fetching an ArticleTag by its ID.
func TestGetArticleTagByID(t *testing.T) {
	article := createRandomArticle(t)
	tag := createRandomTag(t)
	defer cleanupArticle(t, article.ID)
	defer cleanupTag(t, tag.ID)

	created, err := testQueries.CreateArticleTag(context.Background(), CreateArticleTagParams{
		ArticleID: article.ID,
		TagID:     tag.ID,
	})
	require.NoError(t, err)

	fetched, err := testQueries.GetArticleTagByID(context.Background(), created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, fetched.ID)
	require.Equal(t, article.ID, fetched.ArticleID)
	require.Equal(t, tag.ID, fetched.TagID)
}

// TestListTagsByArticle tests listing tags associated with a specific article.
func TestListTagsByArticle(t *testing.T) {
	article := createRandomArticle(t)
	tag1 := createRandomTag(t)
	tag2 := createRandomTag(t)
	defer cleanupArticle(t, article.ID)
	defer cleanupTag(t, tag1.ID)
	defer cleanupTag(t, tag2.ID)

	_, err := testQueries.CreateArticleTag(context.Background(), CreateArticleTagParams{
		ArticleID: article.ID,
		TagID:     tag1.ID,
	})
	require.NoError(t, err)

	tags, err := testQueries.ListTagsByArticle(context.Background(), ListTagsByArticleParams{
		ArticleID: article.ID,
		Limit:     10,
		Offset:    0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, tags)
	require.True(t, tagInList(tags, tag1.ID))
}

// TestListArticlesByTagID tests listing articles associated with a specific tag.
func TestListArticlesByTagID(t *testing.T) {
	article := createRandomArticle(t)
	tag := createRandomTag(t)
	defer cleanupArticle(t, article.ID)
	defer cleanupTag(t, tag.ID)

	_, err := testQueries.CreateArticleTag(context.Background(), CreateArticleTagParams{
		ArticleID: article.ID,
		TagID:     tag.ID,
	})
	require.NoError(t, err)

	articles, err := testQueries.ListArticlesByTagID(context.Background(), ListArticlesByTagIDParams{
		TagID:  tag.ID,
		Limit:  10,
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, articles)
	found := false
	for _, item := range articles {
		if item.ID == article.ID {
			found = true
			break
		}
	}
	require.True(t, found)
}

// TestDeleteArticleTag tests deleting an ArticleTag association.
func TestDeleteArticleTag(t *testing.T) {
	article := createRandomArticle(t)
	tag := createRandomTag(t)
	defer cleanupArticle(t, article.ID)
	defer cleanupTag(t, tag.ID)

	created, err := testQueries.CreateArticleTag(context.Background(), CreateArticleTagParams{
		ArticleID: article.ID,
		TagID:     tag.ID,
	})
	require.NoError(t, err)

	err = testQueries.DeleteArticleTag(context.Background(), DeleteArticleTagParams{
		ArticleID: article.ID,
		TagID:     tag.ID,
	})
	require.NoError(t, err)

	_, err = testQueries.GetArticleTagByID(context.Background(), created.ID)
	require.Error(t, err)
}
