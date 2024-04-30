package ingredient_test

import (
	"connectrpc.com/connect"
	"context"
	_ "github.com/jackc/pgx/v5/stdlib" // "pgx" driver
	"gluttony/internal/database/sqldb"
	"gluttony/internal/ingredient"
	ingredientv1 "gluttony/internal/proto/ingredient/v1"
	"gluttony/internal/x/connectx"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("Create valid ingredient", func(t *testing.T) {
		t.Parallel()
		client := ingredient.NewTestEnv(t, sqldb.NewTestPGXPool(t))

		_, err := client.Create(context.Background(), connect.NewRequest(&ingredientv1.CreateRequest{
			Name:   "Apple",
			Locale: "pl",
		}))
		if err != nil {
			t.Errorf("Expected err to be nil, got %v", err)
		}
	})

	t.Run("Create ingredient unknown locale", func(t *testing.T) {
		t.Parallel()
		client := ingredient.NewTestEnv(t, sqldb.NewTestPGXPool(t))

		_, err := client.Create(context.Background(), connect.NewRequest(&ingredientv1.CreateRequest{
			Name:   "Apple",
			Locale: "de",
		}))

		connectErr, ok := connectx.AsConnectError(err)
		if !ok {
			t.Errorf("Expected err to be *connect.Errror")
			return
		}

		if connectErr.Code() != connect.CodeInvalidArgument {
			t.Errorf("Expected code InvalidArgument, got %v", connectErr.Code())
		}
	})
}

func TestAll(t *testing.T) {
	pool := sqldb.NewTestPGXPool(t)
	_, err := pool.Exec(
		context.Background(),
		`
			INSERT INTO ingredients(id, name) VALUES 
				(1, '{"en": "Apple", "pl": "Jabłko"}'),
				(2, '{"en": "Chicken", "pl": "Kurczak"}'),
				(3, '{"en": "Carrot"}'),
				(4, '{"pl": "Wołowina"}');
		`,
	)
	if err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	t.Run("All en ingredients", func(t *testing.T) {
		t.Parallel()

		client := ingredient.NewTestEnv(t, pool)

		input := connect.NewRequest(&ingredientv1.AllRequest{
			Offset: 0,
			Limit:  20,
			Search: "",
		})
		input.Header().Set("Accept-Language", "en")

		all, err := client.All(context.Background(), input)
		if err != nil {
			t.Errorf("Expected nil error: %v", err)
			return
		}

		if len(all.Msg.Ingredients) != 3 {
			t.Errorf("Expected 3 msg.Ingredients, got %d", len(all.Msg.Ingredients))
			return
		}

		if all.Msg.Ingredients[0].Name != "Apple" {
			t.Errorf("Expected Apple, got %s", all.Msg.Ingredients[0].Name)
		}
	})

	t.Run("All en ingredients filtered", func(t *testing.T) {
		t.Parallel()

		client := ingredient.NewTestEnv(t, pool)

		input := connect.NewRequest(&ingredientv1.AllRequest{
			Offset: 0,
			Limit:  20,
			Search: "Carrot",
		})
		input.Header().Set("Accept-Language", "en")

		all, err := client.All(context.Background(), input)
		if err != nil {
			t.Errorf("Expected nil error: %v", err)
			return
		}

		if len(all.Msg.Ingredients) != 1 {
			t.Errorf("Expected 1 msg.Ingredients, got %d", len(all.Msg.Ingredients))
		}
	})

	t.Run("All en ingredients paginated", func(t *testing.T) {
		t.Parallel()

		client := ingredient.NewTestEnv(t, pool)

		input := connect.NewRequest(&ingredientv1.AllRequest{
			Offset: 0,
			Limit:  2,
			Search: "",
		})
		input.Header().Set("Accept-Language", "en")

		all, err := client.All(context.Background(), input)
		if err != nil {
			t.Errorf("Expected nil error: %v", err)
			return
		}

		if len(all.Msg.Ingredients) != 2 {
			t.Errorf("Expected 2 msg.Ingredients, got %d", len(all.Msg.Ingredients))
		}

		if all.Msg.Ingredients[0].Id != 1 {
			t.Errorf("Expected id 1, got %d", all.Msg.Ingredients[0].Id)
		}

		if all.Msg.Ingredients[0].Id != 1 {
			t.Errorf("Expected id 2, got %d", all.Msg.Ingredients[1].Id)
		}

		input = connect.NewRequest(&ingredientv1.AllRequest{
			Offset: 2,
			Limit:  2,
			Search: "",
		})
		input.Header().Set("Accept-Language", "en")

		all, err = client.All(context.Background(), input)
		if err != nil {
			t.Errorf("Expected nil error: %v", err)
			return
		}

		if len(all.Msg.Ingredients) != 1 {
			t.Errorf("Expected 2 msg.Ingredients, got %d", len(all.Msg.Ingredients))
		}

		if all.Msg.Ingredients[0].Id != 3 {
			t.Errorf("Expected id 1, got %d", all.Msg.Ingredients[0].Id)
		}
	})
}
