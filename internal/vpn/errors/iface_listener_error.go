package errors

import "bvpn-prototype/internal/infrastructure/errors"

func InterfaceListenerError(data ...any) errors.Error {
	return errors.Error{
		Code: 11311,
		Type: "InterfaceListenerError",
		Data: data,
	}
}
