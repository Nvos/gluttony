package user

import (
	"context"
	"errors"
	"fmt"
	"gluttony/internal/auth"
	v1 "gluttony/internal/proto/user/v1"
	"gluttony/internal/proto/user/v1/userv1connect"
	"net/http"

	"connectrpc.com/connect"
)

var _ userv1connect.UserServiceHandler = (*ConnectService)(nil)

type ConnectService struct {
	service *Service
}

func (c *ConnectService) Logout(
	ctx context.Context,
	_ *connect.Request[v1.LogoutRequest],
) (*connect.Response[v1.LogoutResponse], error) {
	token, err := auth.GetSessionToken(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	if err := c.service.Logout(ctx, token); err != nil {
		return nil, fmt.Errorf("logout: %w", err)
	}

	cookie := auth.NewExpiredSessionCookie()
	out := connect.NewResponse(&v1.LogoutResponse{})
	out.Header().Set("Set-Cookie", cookie.String())

	return out, nil
}

func (c *ConnectService) Login(
	ctx context.Context,
	r *connect.Request[v1.LoginRequest],
) (*connect.Response[v1.LoginResponse], error) {
	input := LoginInput{
		Username: r.Msg.Username,
		Password: r.Msg.Password,
	}

	token, err := c.service.Login(ctx, input)
	if err != nil && errors.Is(err, ErrInvalidCredentials) {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	if err != nil {
		return nil, err
	}

	cookie := auth.NewUnsecureSessionCookie(token)
	out := &v1.LoginResponse{
		User: &v1.User{
			Username: r.Msg.Username,
		},
	}

	resp := connect.NewResponse(out)
	resp.Header().Set("Set-Cookie", cookie.String())
	resp.Header().Set("Cache-Control", `no-cache="Set-Cookie"`)

	return resp, nil
}

func (c *ConnectService) Me(
	ctx context.Context,
	_ *connect.Request[v1.MeRequest],
) (*connect.Response[v1.MeResponse], error) {
	session, err := auth.GetSession[Session](ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	out := &v1.MeResponse{
		User: &v1.User{
			Username: session.Username,
		},
	}

	return connect.NewResponse(out), nil
}

func NewConnectHandler(service *Service, opts ...connect.HandlerOption) (string, http.Handler, error) {
	if service == nil {
		return "", nil, fmt.Errorf("new connect handler: nil service")
	}

	connectService := &ConnectService{
		service: service,
	}

	path, handler := userv1connect.NewUserServiceHandler(connectService, opts...)
	return path, handler, nil
}
