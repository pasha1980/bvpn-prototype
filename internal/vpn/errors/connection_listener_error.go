package errors

import "bvpn-prototype/internal/infrastructure/errors"

func ConnectionListenerError(data ...any) errors.Error {
	return errors.Error{
		Code: 11311,
		Type: "ConnectionListenerError",
		Data: data,
	}
}
