package recipe

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"gluttony/internal/database/pagination"
	v1 "gluttony/internal/proto/recipe/v1"
	"gluttony/internal/proto/recipe/v1/recipev1connect"
	"net/http"
)

var _ recipev1connect.RecipeServiceHandler = (*ConnectService)(nil)

// TODO(AK) 06/03/2024: proper errors, need to know which one is validation

type ConnectService struct {
	store Store
}

func (s *ConnectService) SingleRecipe(
	ctx context.Context,
	r *connect.Request[v1.SingleRecipeRequest],
) (*connect.Response[v1.SingleRecipeResponse], error) {

	single, err := s.store.Single(ctx, r.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	all, err := s.store.All(ctx, r.Msg.Search, offsetPagination)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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

	steps := make([]CreateStep, 0, len(r.Msg.Steps))
	for i := range r.Msg.Steps {
		v := r.Msg.Steps[i]
		step := CreateStep{
			Order:       v.Order,
			Description: v.Description,
		}

		if err := ValidateCreateRecipeStep(step); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		steps = append(steps, step)
	}

	in := CreateRecipe{
		Name:        r.Msg.Name,
		Description: r.Msg.Description,
		Steps:       steps,
	}

	if err := ValidateCreateRecipe(in); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	recipeID, err := s.store.Create(ctx, in)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.CreateRecipeResponse{
		RecipeId: recipeID,
	}), nil
}

func NewConnectHandler(store Store) (string, http.Handler, error) {
	if store == nil {
		return "", nil, fmt.Errorf("new connect handler: nil store")
	}

	service := &ConnectService{
		store: store,
	}

	path, handler := recipev1connect.NewRecipeServiceHandler(service)
	return path, handler, nil
}
