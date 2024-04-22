package recipe

import (
	"context"
	"fmt"
	"gluttony/internal/i18n"
	"gluttony/internal/x/validatex"
	"net/http"

	"connectrpc.com/connect"

	v1 "gluttony/internal/proto/recipe/v1"
	"gluttony/internal/proto/recipe/v1/recipev1connect"
)

var _ recipev1connect.RecipeServiceHandler = (*ConnectService)(nil)

type ConnectService struct {
	service *Service
}

func (s *ConnectService) SingleRecipe(
	ctx context.Context,
	r *connect.Request[v1.SingleRecipeRequest],
) (*connect.Response[v1.SingleRecipeResponse], error) {
	locale, err := i18n.GetLocale(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	single, err := s.service.SingleRecipe(ctx, locale, r.Msg.Id)
	if err != nil {
		return nil, err
	}

	out := &v1.SingleRecipeResponse{
		Id:   single.ID,
		Name: single.Name,
		//Content: single.
	}

	return connect.NewResponse(out), nil
}

func (s *ConnectService) AllRecipes(
	ctx context.Context,
	r *connect.Request[v1.AllRecipesRequest],
) (*connect.Response[v1.AllRecipesResponse], error) {
	locale, err := i18n.GetLocale(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	input, err := NewAllRecipesInput(locale, r.Msg)

	all, err := s.service.AllRecipes(ctx, input)
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
		return nil, validatex.SerializeAsConnect(err)
	}

	// TODO, 20/03/2024: return recipe instead
	if _, err := s.service.CreateRecipe(ctx, create); err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.CreateRecipeResponse{}), nil
}

func NewConnectHandler(
	service *Service,
	opts ...connect.HandlerOption,
) (string, http.Handler, error) {
	if service == nil {
		return "", nil, fmt.Errorf("new connect handler: nil service")
	}

	cs := &ConnectService{
		service: service,
	}

	path, handler := recipev1connect.NewRecipeServiceHandler(cs, opts...)
	return path, handler, nil
}
