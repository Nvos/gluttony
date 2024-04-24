package ingredient

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"gluttony/internal/i18n"
	v1 "gluttony/internal/proto/ingredient/v1"
	"gluttony/internal/proto/ingredient/v1/ingredientv1connect"
	"net/http"
)

var _ ingredientv1connect.IngredientServiceHandler = (*ConnectServiceV1)(nil)

type ConnectServiceV1 struct {
	service *Service
}

func (c *ConnectServiceV1) All(
	ctx context.Context,
	r *connect.Request[v1.AllRequest],
) (*connect.Response[v1.AllResponse], error) {
	locale, err := i18n.GetLocale(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	input, err := NewAllIngredientsInput(locale, r.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	all, err := c.service.All(ctx, input)
	if err != nil {
		return nil, err
	}

	out := make([]*v1.Ingredient, 0, len(all))
	for i := range all {
		out = append(out, &v1.Ingredient{
			Id:   all[i].ID,
			Name: all[i].Name,
		})
	}

	return connect.NewResponse(&v1.AllResponse{Ingredients: out}), nil
}

func (c *ConnectServiceV1) Create(
	ctx context.Context,
	r *connect.Request[v1.CreateRequest],
) (*connect.Response[v1.CreateResponse], error) {
	input, err := NewCreateIngredientInput(r.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	if err := c.service.Create(ctx, input); err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.CreateResponse{}), nil
}

func NewConnectHandler(
	service *Service,
	opts ...connect.HandlerOption,
) (string, http.Handler, error) {
	if service == nil {
		return "", nil, fmt.Errorf("new connect handler: nil service")
	}

	cs := &ConnectServiceV1{
		service: service,
	}

	path, handler := ingredientv1connect.NewIngredientServiceHandler(cs, opts...)
	return path, handler, nil
}

func NewConnectClient(client connect.HTTPClient, baseURL string) ingredientv1connect.IngredientServiceClient {
	return ingredientv1connect.NewIngredientServiceClient(client, baseURL)
}
