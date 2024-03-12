package validateutil

import (
	"connectrpc.com/connect"
	"errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func SerializeAsConnect(err error) error {
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		return nil
	}

	if len(validationErr.Violations) == 0 {
		return nil
	}

	violations := make([]*errdetails.BadRequest_FieldViolation, 0, len(validationErr.Violations))
	for i := range validationErr.Violations {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       validationErr.Violations[i].Path,
			Description: validationErr.Violations[i].Rule,
		})
	}

	badRequest := &errdetails.BadRequest{
		FieldViolations: violations,
	}

	connectErr := connect.NewError(connect.CodeInvalidArgument, errors.New("validation failed"))
	if detail, detailErr := connect.NewErrorDetail(badRequest); detailErr == nil {
		connectErr.AddDetail(detail)

		return connectErr
	}

	return errors.New("unable to serialize validation err as connect err")
}
