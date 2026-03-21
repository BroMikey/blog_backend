package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

// createRandomArticle creates a test article and returns it for test assertions
func createRandomArticle(t *testing.T) Article {
	t.Helper()

	args := CreateArticleParams{
		Title:     "Test Article",
		Content:   "This is a test article.",
		AuthorUid: sql.NullInt64{Int64: 1, Valid: true},
	}
	article, err := testQueries.CreateArticle(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, article)

	return article
}

// cleanupArticle deletes an article used in tests
func cleanupArticle(t *testing.T, articleID int64) {
	t.Helper()

	err := testQueries.DeleteArticle(context.Background(), articleID)
	require.NoError(t, err)
}

func TestCreateArticle(t *testing.T) {
	article := createRandomArticle(t)

	args := CreateArticleParams{
		Title:     "Test Article",
		Content:   "This is a test article.",
		AuthorUid: sql.NullInt64{Int64: 1, Valid: true},
	}

	require.Equal(t, args.Title, article.Title)
	require.Equal(t, args.Content, article.Content)
	require.Equal(t, args.AuthorUid, article.AuthorUid)

	require.NotZero(t, article.ID)
	require.Equal(t, int32(0), article.LikeCount)
	require.Equal(t, int32(0), article.CommentCount)
}

func TestListArticles(t *testing.T) {
	createRandomArticle(t)

	listArgs := ListArticlesParams{
		Limit:  5,
		Offset: 0,
	}

	articles, err := testQueries.ListArticles(context.Background(), listArgs)
	require.NoError(t, err)
	require.NotEmpty(t, articles)
}

func TestGetArticle(t *testing.T) {
	article := createRandomArticle(t)

	fetchedArticle, err := testQueries.GetArticle(context.Background(), article.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedArticle)

	require.Equal(t, article.ID, fetchedArticle.ID)
	require.Equal(t, article.Title, fetchedArticle.Title)
	require.Equal(t, article.Content, fetchedArticle.Content)
	require.Equal(t, article.AuthorUid, fetchedArticle.AuthorUid)
}

func TestDecrementArticleCommentCount(t *testing.T) {
	article := createRandomArticle(t)

	// Increment comment count first
	commentcount, err := testQueries.IncrementArticleCommentCount(context.Background(), article.ID)
	require.NoError(t, err)
	require.Equal(t, commentcount, int32(1))

	// Now decrement
	commentcount, err = testQueries.DecrementArticleCommentCount(context.Background(), article.ID)
	require.NoError(t, err)
	require.Equal(t, commentcount, int32(0))

	updatedArticle, err := testQueries.GetArticle(context.Background(), article.ID)
	require.NoError(t, err)
	require.Equal(t, int32(0), updatedArticle.CommentCount)
}

func TestDecrementArticleLikeCount(t *testing.T) {
	article := createRandomArticle(t)

	// Increment like count first
	likecount, err := testQueries.IncrementArticleLikeCount(context.Background(), article.ID)
	require.NoError(t, err)
	require.Equal(t, likecount, int32(1))

	// Now decrement
	likecount, err = testQueries.DecrementArticleLikeCount(context.Background(), article.ID)
	require.NoError(t, err)
	require.Equal(t, likecount, int32(0))

	updatedArticle, err := testQueries.GetArticle(context.Background(), article.ID)
	require.NoError(t, err)
	require.Equal(t, int32(0), updatedArticle.LikeCount)
}

func TestIncrementArticleCommentCount(t *testing.T) {
	article := createRandomArticle(t)

	// Increment comment count
	commentcount, err := testQueries.IncrementArticleCommentCount(context.Background(), article.ID)
	require.NoError(t, err)
	require.Equal(t, commentcount, int32(1))

	updatedArticle, err := testQueries.GetArticle(context.Background(), article.ID)
	require.NoError(t, err)
	require.Equal(t, int32(1), updatedArticle.CommentCount)
}

func TestIncrementArticleLikeCount(t *testing.T) {
	article := createRandomArticle(t)

	// Increment like count
	likecount, err := testQueries.IncrementArticleLikeCount(context.Background(), article.ID)
	require.NoError(t, err)
	require.Equal(t, likecount, int32(1))

	// update article and check
	updatedArticle, err := testQueries.GetArticle(context.Background(), article.ID)
	require.NoError(t, err)
	require.Equal(t, int32(1), updatedArticle.LikeCount)
}

func TestIncrementArticleViewCount(t *testing.T) {
	article := createRandomArticle(t)

	// Increment view count
	viewcount, err := testQueries.IncrementArticleViewCount(context.Background(), article.ID)
	require.NoError(t, err)
	require.Equal(t, viewcount, int32(1))

	updatedArticle, err := testQueries.GetArticle(context.Background(), article.ID)
	require.NoError(t, err)
	require.Equal(t, int32(1), updatedArticle.ViewCount)
}

// test listing articles by category
func TestListArticleByCategory(t *testing.T) {
	article := createRandomArticle(t)
	category := createRandomCategory(t)

	testQueries.CreateArticleCategory(context.Background(), CreateArticleCategoryParams{
		ArticleID:  article.ID,
		CategoryID: category.ID,
	})
	listarticle, err := testQueries.ListArticleByCategory(context.Background(), ListArticleByCategoryParams{
		CategoryID: category.ID,
		Limit:      10,
		Offset:     0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, listarticle)

	for _, a := range listarticle {
		if a.ID == article.ID {
			return
		}
	}

	require.Fail(t, "article not found in list by category")
}

// test listing articles by tag
func TestListArticleByTag(t *testing.T) {
	article := createRandomArticle(t)
	tag := createRandomTag(t)

	testQueries.CreateArticleTag(context.Background(), CreateArticleTagParams{
		ArticleID: article.ID,
		TagID:     tag.ID,
	})
	listarticle, err := testQueries.ListArticleByTag(context.Background(), ListArticleByTagParams{
		TagID:  tag.ID,
		Limit:  10,
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, listarticle)

	for _, a := range listarticle {
		if a.ID == article.ID {
			return
		}
	}

	require.Fail(t, "article not found in list by tag")
}

// TestUpdateArticle tests updating an article's title and content.
func TestUpdateArticle(t *testing.T) {
	article := createRandomArticle(t)

	updateArgs := UpdateArticleParams{
		Title:   "Updated Title",
		Content: "Updated content.",
		ID:      article.ID,
	}
	updatedArticle, err := testQueries.UpdateArticle(context.Background(), updateArgs)
	require.NoError(t, err)
	require.NotEmpty(t, updatedArticle)

	require.Equal(t, updateArgs.Title, updatedArticle.Title)
	require.Equal(t, updateArgs.Content, updatedArticle.Content)
}

func TestListArticlesByCommentCount(t *testing.T) {
	article := createRandomArticle(t)
	for i := 0; i < 5; i++ {
		createRandomComment(t, article.ID)
	}

	listComment, err := testQueries.ListArticlesByCommentCount(context.Background(), ListArticlesByCommentCountParams{
		Limit:  100000,
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, listComment)

	for _, a := range listComment {
		if a.ID == article.ID {
			return
		}
	}

	require.Fail(t, "article not found in list by comment count")
}

func TestListArticlesByLikeCount(t *testing.T) {
	article := createRandomArticle(t)

	for i := 0; i < 5; i++ {
		user := createRandomUser(t)
		like, err := testQueries.CreateArticleLike(context.Background(), CreateArticleLikeParams{
			Uid:       user.Uid,
			ArticleID: article.ID,
		})
		require.NoError(t, err)
		require.NotEmpty(t, like)
	}
	listLike, err := testQueries.ListArticlesByLikeCount(context.Background(), ListArticlesByLikeCountParams{
		Limit:  100000,
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, listLike)

	for _, a := range listLike {
		if a.ID == article.ID {
			return
		}
	}

	require.Fail(t, "article not found in list by like count")
}

func TestListArticlesByViewCount(t *testing.T) {
	article := createRandomArticle(t)
	for i := 0; i < 5; i++ {
		_, err := testQueries.IncrementArticleViewCount(context.Background(), article.ID)
		require.NoError(t, err)
	}

	listView, err := testQueries.ListArticlesByViewCount(context.Background(), ListArticlesByViewCountParams{
		Limit:  100000,
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, listView)

	for _, a := range listView {
		if a.ID == article.ID {
			return
		}
	}

	require.Fail(t, "article not found in list by view count")
}
