package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateArticleCategory(t *testing.T) {
	article := createRandomArticle(t)
	category := createRandomCategory(t)

	created, err := testQueries.CreateArticleCategory(context.Background(), CreateArticleCategoryParams{ArticleID: article.ID, CategoryID: category.ID})
	require.NoError(t, err)
	require.Equal(t, article.ID, created.ArticleID)
	require.Equal(t, category.ID, created.CategoryID)
}

func TestGetArticleCategoryByID(t *testing.T) {
	article := createRandomArticle(t)
	category := createRandomCategory(t)

	created, err := testQueries.CreateArticleCategory(context.Background(), CreateArticleCategoryParams{ArticleID: article.ID, CategoryID: category.ID})
	require.NoError(t, err)

	fetched, err := testQueries.GetArticleCategoryByID(context.Background(), created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, fetched.ID)
	require.Equal(t, article.ID, fetched.ArticleID)
	require.Equal(t, category.ID, fetched.CategoryID)
}

func TestListCategoriesByArticle(t *testing.T) {
	article := createRandomArticle(t)
	category := createRandomCategory(t)

	_, err := testQueries.CreateArticleCategory(context.Background(), CreateArticleCategoryParams{ArticleID: article.ID, CategoryID: category.ID})
	require.NoError(t, err)

	categories, err := testQueries.ListCategoriesByArticle(context.Background(), ListCategoriesByArticleParams{ArticleID: article.ID, Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.NotEmpty(t, categories)
	require.True(t, categoryInList(categories, category.ID))
}

func TestListArticlesByCategoryID(t *testing.T) {
	article := createRandomArticle(t)
	category := createRandomCategory(t)

	_, err := testQueries.CreateArticleCategory(context.Background(), CreateArticleCategoryParams{ArticleID: article.ID, CategoryID: category.ID})
	require.NoError(t, err)

	articles, err := testQueries.ListArticlesByCategoryID(context.Background(), ListArticlesByCategoryIDParams{CategoryID: category.ID, Limit: 10, Offset: 0})
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

func TestDeleteArticleCategory(t *testing.T) {
	article := createRandomArticle(t)
	category := createRandomCategory(t)

	created, err := testQueries.CreateArticleCategory(context.Background(), CreateArticleCategoryParams{ArticleID: article.ID, CategoryID: category.ID})
	require.NoError(t, err)

	err = testQueries.DeleteArticleCategory(context.Background(), DeleteArticleCategoryParams{ArticleID: article.ID, CategoryID: category.ID})
	require.NoError(t, err)

	_, err = testQueries.GetArticleCategoryByID(context.Background(), created.ID)
	require.Error(t, err)
}
