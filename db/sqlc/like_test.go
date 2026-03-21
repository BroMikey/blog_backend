package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// likeInList checks if a like by a specific user on a specific article exists in the list.
func likeInList(items []ArticleLike, uid, articleID int64) bool {
	for _, item := range items {
		if item.Uid == uid && item.ArticleID == articleID {
			return true
		}
	}
	return false
}

// TestCreateArticleLike tests creating a like on an article.
func TestCreateArticleLike(t *testing.T) {
	user := createRandomUser(t)
	article := createRandomArticle(t)

	created, err := testQueries.CreateArticleLike(context.Background(), CreateArticleLikeParams{Uid: user.Uid, ArticleID: article.ID})
	require.NoError(t, err)
	require.Equal(t, user.Uid, created.Uid)
	require.Equal(t, article.ID, created.ArticleID)
}

// TestCreateArticleLike tests creating a like on an article.
func TestGetArticleLike(t *testing.T) {
	user := createRandomUser(t)
	article := createRandomArticle(t)

	created, err := testQueries.CreateArticleLike(context.Background(), CreateArticleLikeParams{Uid: user.Uid, ArticleID: article.ID})
	require.NoError(t, err)

	fetched, err := testQueries.GetArticleLike(context.Background(), GetArticleLikeParams{Uid: user.Uid, ArticleID: article.ID})
	require.NoError(t, err)
	require.Equal(t, created.ID, fetched.ID)
}

// TestGetArticleLikeByID tests fetching an ArticleLike by its ID.
func TestGetArticleLikeByID(t *testing.T) {
	user := createRandomUser(t)
	article := createRandomArticle(t)

	created, err := testQueries.CreateArticleLike(context.Background(), CreateArticleLikeParams{Uid: user.Uid, ArticleID: article.ID})
	require.NoError(t, err)

	fetched, err := testQueries.GetArticleLikeByID(context.Background(), created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, fetched.ID)
}

// TestListArticleLikesByArticle tests listing likes associated with a specific article.
func TestListArticleLikesByArticle(t *testing.T) {
	user := createRandomUser(t)
	article := createRandomArticle(t)

	_, err := testQueries.CreateArticleLike(context.Background(), CreateArticleLikeParams{Uid: user.Uid, ArticleID: article.ID})
	require.NoError(t, err)

	likes, err := testQueries.ListArticleLikesByArticle(context.Background(), ListArticleLikesByArticleParams{ArticleID: article.ID, Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.True(t, likeInList(likes, user.Uid, article.ID))
}

// TestDeleteArticleLike tests deleting an ArticleLike association.
func TestDeleteArticleLike(t *testing.T) {
	user := createRandomUser(t)
	article := createRandomArticle(t)

	created, err := testQueries.CreateArticleLike(context.Background(), CreateArticleLikeParams{Uid: user.Uid, ArticleID: article.ID})
	require.NoError(t, err)

	err = testQueries.DeleteArticleLike(context.Background(), DeleteArticleLikeParams{Uid: user.Uid, ArticleID: article.ID})
	require.NoError(t, err)

	_, err = testQueries.GetArticleLikeByID(context.Background(), created.ID)
	require.Error(t, err)
}
