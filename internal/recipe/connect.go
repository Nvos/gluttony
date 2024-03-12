package recipe

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"gluttony/internal/database/pagination"
	v1 "gluttony/internal/proto/recipe/v1"
	"gluttony/internal/proto/recipe/v1/recipev1connect"
	"gluttony/internal/util/validateutil"
	"net/http"
)

var _ recipev1connect.RecipeServiceHandler = (*ConnectService)(nil)

type ConnectService struct {
	store Store
}

func (s *ConnectService) SingleRecipe(
	ctx context.Context,
	r *connect.Request[v1.SingleRecipeRequest],
) (*connect.Response[v1.SingleRecipeResponse], error) {

	single, err := s.store.Single(ctx, r.Msg.Id)
	if err != nil {
		return nil, err
	}

	steps := make([]*v1.RecipeStep, 0, len(single.Steps))
	for i := range single.Steps {
		v := single.Steps[i]
		steps = append(steps, &v1.RecipeStep{
			Id:          v.ID,
			Order:       v.Order,
			Description: v.Description,
		})
	}

	out := &v1.SingleRecipeResponse{
		Id:    single.ID,
		Name:  single.Name,
		Steps: steps,
	}

	return connect.NewResponse(out), nil
}

func (s *ConnectService) AllRecipes(
	ctx context.Context,
	r *connect.Request[v1.AllRecipesRequest],
) (*connect.Response[v1.AllRecipesResponse], error) {

	offsetPagination := pagination.OffsetPagination{
		Offset: r.Msg.Offset,
		Limit:  r.Msg.Limit,
	}

	if err := pagination.ValidateOffsetPagination(offsetPagination); err != nil {
		return nil, err
	}

	all, err := s.store.All(ctx, r.Msg.Search, offsetPagination)
	if err != nil {
		return nil, err
	}

	recipes := make([]*v1.Recipe, 0, len(all))
	for i := range all {
		v := all[i]

		recipes = append(recipes, &v1.Recipe{
			Id:          v.ID,
			Name:        v.Name,
			Description: v.Description,
		})
	}

	return connect.NewResponse(&v1.AllRecipesResponse{
		Recipes: recipes,
	}), nil
}

func (s *ConnectService) CreateRecipe(
	ctx context.Context,
	r *connect.Request[v1.CreateRecipeRequest],
) (*connect.Response[v1.CreateRecipeResponse], error) {

	create, err := NewCreateRecipe(r.Msg)
	if err != nil {
		return nil, validateutil.SerializeAsConnect(err)
	}

	recipeID, err := s.store.Create(ctx, create)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.CreateRecipeResponse{
		RecipeId: recipeID,
	}), nil
}

func NewConnectHandler(store Store, opts ...connect.HandlerOption) (string, http.Handler, error) {
	if store == nil {
		return "", nil, fmt.Errorf("new connect handler: nil store")
	}

	service := &ConnectService{
		store: store,
	}

	path, handler := recipev1connect.NewRecipeServiceHandler(service, opts...)
	return path, handler, nil
}
