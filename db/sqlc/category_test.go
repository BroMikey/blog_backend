package db

import (
	"context"
	"testing"

	"github.com/BroMikey/blog_backend/utils"
	"github.com/stretchr/testify/require"
)

func createRandomCategory(t *testing.T) Category {
	t.Helper()

	args := CreateCategoryParams{
		Name: utils.RandomString(8),
		Sort: int16(utils.RandomInt(0, 10)),
	}
	category, err := testQueries.CreateCategory(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, category)
	require.Equal(t, args.Name, category.Name)
	require.Equal(t, args.Sort, category.Sort)
	require.NotZero(t, category.ID)
	require.NotZero(t, category.CreatedAt)

	return category
}

func cleanupCategory(t *testing.T, id int64) {
	t.Helper()

	err := testQueries.DeleteCategory(context.Background(), id)
	require.NoError(t, err)
}

func categoryInList(categories []Category, id int64) bool {
	for _, category := range categories {
		if category.ID == id {
			return true
		}
	}
	return false
}

func TestCreateCategory(t *testing.T) {
	category := createRandomCategory(t)

	require.NotZero(t, category.ID)
}

func TestGetCategoryByID(t *testing.T) {
	category := createRandomCategory(t)

	fetched, err := testQueries.GetCategoryByID(context.Background(), category.ID)
	require.NoError(t, err)
	require.Equal(t, category.ID, fetched.ID)
	require.Equal(t, category.Name, fetched.Name)
	require.Equal(t, category.Sort, fetched.Sort)
}

func TestListCategories(t *testing.T) {
	category1 := createRandomCategory(t)
	category2 := createRandomCategory(t)

	categories, err := testQueries.ListCategories(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, categories)
	require.True(t, categoryInList(categories, category1.ID))
	require.True(t, categoryInList(categories, category2.ID))
}

func TestUpdateCategory(t *testing.T) {
	category := createRandomCategory(t)

	updated, err := testQueries.UpdateCategory(context.Background(), UpdateCategoryParams{
		ID:   category.ID,
		Name: utils.RandomString(10),
		Sort: int16(utils.RandomInt(0, 10)),
	})
	require.NoError(t, err)
	require.Equal(t, category.ID, updated.ID)
}

func TestDeleteCategory(t *testing.T) {
	category := createRandomCategory(t)
	err := testQueries.DeleteCategory(context.Background(), category.ID)
	require.NoError(t, err)

	_, err = testQueries.GetCategoryByID(context.Background(), category.ID)
	require.Error(t, err)
}
