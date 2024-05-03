package ingredient_test

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"gluttony/internal/database/pagination"
	"gluttony/internal/i18n"
	"gluttony/internal/ingredient"
	"gluttony/internal/x/testing/assert"
	"gluttony/internal/x/testing/testdb"
	"testing"
)

func newService(t *testing.T, pool *pgxpool.Pool) *ingredient.Service {
	t.Helper()

	store := ingredient.NewStorePostgres(pool)
	service := ingredient.NewService(store)

	return service
}

func TestService_Create(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		pool := testdb.NewTestPGXPool(t)
		service := newService(t, pool)

		input := ingredient.CreateIngredientInput{
			Locale: i18n.EnLocale,
			Name:   "Pepper",
			Unit:   ingredient.Weight,
		}

		if err := service.Create(context.Background(), input); assert.NilErr(t, err) {
			t.FailNow()
		}
	})
}

func TestService_All(t *testing.T) {
	tests := []struct {
		name     string
		input    ingredient.AllIngredientsInput
		expected []ingredient.Ingredient
	}{
		{
			name: "All en ingredients",
			input: ingredient.AllIngredientsInput{
				Locale:     i18n.EnLocale,
				Pagination: pagination.OffsetPagination{Limit: 999},
			},
			expected: []ingredient.Ingredient{
				{ID: 1, Name: "Apple"},
				{ID: 2, Name: "Chicken"},
				{ID: 3, Name: "Carrot"},
			},
		},
		{
			name: "All pl ingredients",
			input: ingredient.AllIngredientsInput{
				Locale:     i18n.PlLocale,
				Pagination: pagination.OffsetPagination{Limit: 999},
			},
			expected: []ingredient.Ingredient{
				{ID: 1, Name: "Jabłko"},
				{ID: 2, Name: "Kurczak"},
				{ID: 4, Name: "Wołowina"},
				{ID: 5, Name: "Śliwka"},
			},
		},
		{
			name: "Pagination offset 0",
			input: ingredient.AllIngredientsInput{
				Locale:     i18n.EnLocale,
				Pagination: pagination.OffsetPagination{Offset: 0, Limit: 2},
			},
			expected: []ingredient.Ingredient{
				{ID: 1, Name: "Apple"},
				{ID: 2, Name: "Chicken"},
			},
		},
		{
			name: "Pagination offset 2",
			input: ingredient.AllIngredientsInput{
				Locale:     i18n.EnLocale,
				Pagination: pagination.OffsetPagination{Offset: 2, Limit: 2},
			},
			expected: []ingredient.Ingredient{
				{ID: 3, Name: "Carrot"},
			},
		},
		{
			name: "Search exact",
			input: ingredient.AllIngredientsInput{
				Locale:     i18n.EnLocale,
				Search:     "Carrot",
				Pagination: pagination.OffsetPagination{Offset: 0, Limit: 20},
			},
			expected: []ingredient.Ingredient{
				{ID: 3, Name: "Carrot"},
			},
		},
		{
			name: "Search similar",
			input: ingredient.AllIngredientsInput{
				Locale:     i18n.EnLocale,
				Search:     "Carro",
				Pagination: pagination.OffsetPagination{Offset: 0, Limit: 20},
			},
			expected: []ingredient.Ingredient{
				{ID: 3, Name: "Carrot"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			pool := testdb.NewTestPGXPool(t)
			testdb.Seed(t, pool, "ingredients.sql")

			service := newService(t, pool)
			got, err := service.All(context.Background(), test.input)
			if assert.NilErr(t, err) {
				t.FailNow()
			}

			if assert.Equal(t, len(got), len(test.expected)) {
				t.FailNow()
			}

			for i := range test.expected {
				assert.Equal(t, got[i].ID, test.expected[i].ID)
				assert.Equal(t, got[i].Name, test.expected[i].Name)
			}
		})
	}
}
