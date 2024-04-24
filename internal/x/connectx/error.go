package connectx

import (
	"connectrpc.com/connect"
	"errors"
)

func AsConnectError(err error) (*connect.Error, bool) {
	var connectErr *connect.Error
	if errors.As(err, &connectErr) {
		return connectErr, true
	}

	return nil, false
}
