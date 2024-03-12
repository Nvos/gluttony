package connectutil

import (
	"connectrpc.com/connect"
	"context"
	"log/slog"
)

func ErrorInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			response, err := next(ctx, req)
			if err != nil {
				if connect.IsWireError(err) {
					return response, err
				}

				logger.Error("Connect handler", slog.String("err", err.Error()))

				return response, connect.NewError(connect.CodeInternal, nil)
			}

			return response, err
		}
	}
}
