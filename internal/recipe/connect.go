package recipe

import (
	"connectrpc.com/connect"
	"context"
	v1 "gluttony/internal/proto/recipe/v1"
	"gluttony/internal/proto/recipe/v1/recipev1connect"
)

var _ recipev1connect.RecipeServiceHandler = (*ConnectService)(nil)

type ConnectService struct {
}

func (c *ConnectService) GetRecipe(
	ctx context.Context,
	request *connect.Request[v1.GetRecipeRequest],
) (*connect.Response[v1.GetRecipeResponse], error) {
	//TODO implement me
	panic("implement me")
}
