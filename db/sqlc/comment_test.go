package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomComment(t *testing.T, id int64) ArticleComment {
	t.Helper()

	args := CreateCommentParams{
		ArticleID: id,
		Uid:       1,
		Content:   "This is a test comment.",
		ParentID:  sql.NullInt64{Int64: 0, Valid: false},
		Status:    1,
	}
	comment, err := testQueries.CreateComment(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, comment)

	return comment
}

func cleanupComment(t *testing.T, commentID int64) {
	t.Helper()

	err := testQueries.DeleteComment(context.Background(), commentID)
	require.NoError(t, err)

}

func TestCreateComment(t *testing.T) {
	article := createRandomArticle(t)

	comment := createRandomComment(t, article.ID)

	args := CreateCommentParams{
		ArticleID: article.ID,
		Uid:       1,
		Content:   "This is a test comment.",
		ParentID:  sql.NullInt64{Int64: 0, Valid: false},
		Status:    1,
	}
	require.Equal(t, args.ArticleID, comment.ArticleID)
	require.Equal(t, args.Uid, comment.Uid)
	require.Equal(t, args.Content, comment.Content)
	require.Equal(t, args.ParentID, comment.ParentID)
	require.Equal(t, args.Status, comment.Status)

	require.NotZero(t, comment.ID)
	require.NotZero(t, comment.CreatedAt)
}

func TestDeleteComment(t *testing.T) {
	article := createRandomArticle(t)

	comment := createRandomComment(t, article.ID)

	err := testQueries.DeleteComment(context.Background(), comment.ID)
	require.NoError(t, err)

	_, err = testQueries.GetComment(context.Background(), comment.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestGetComment(t *testing.T) {
	article := createRandomArticle(t)

	comment := createRandomComment(t, article.ID)

	fetchedComment, err := testQueries.GetComment(context.Background(), comment.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedComment)
	require.Equal(t, comment.ID, fetchedComment.ID)
	require.Equal(t, comment.ArticleID, fetchedComment.ArticleID)
	require.Equal(t, comment.Uid, fetchedComment.Uid)
	require.Equal(t, comment.Content, fetchedComment.Content)
	require.Equal(t, comment.ParentID, fetchedComment.ParentID)
	require.Equal(t, comment.Status, fetchedComment.Status)
}

func TestListComments(t *testing.T) {
	article := createRandomArticle(t)

	// Create multiple test comments for pagination testing
	for i := 0; i < 10; i++ {
		createRandomComment(t, article.ID)
	}

	listArgs := ListCommentsParams{
		ArticleID: article.ID,
		Limit:     5,
		Offset:    0,
	}

	comments, err := testQueries.ListComments(context.Background(), listArgs)
	require.NoError(t, err)
	require.NotEmpty(t, comments)
	require.Len(t, comments, int(listArgs.Limit))

	for _, comment := range comments {
		require.NotEmpty(t, comment)
		require.Equal(t, article.ID, comment.ArticleID)
	}
}
