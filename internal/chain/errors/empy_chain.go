package errors

import "bvpn-prototype/internal/infrastructure/errors"

func EmptyChainError(data ...any) errors.Error {
	return errors.Error{
		Code: 11111,
		Type: "EmptyChainError",
		Data: data,
	}
}
