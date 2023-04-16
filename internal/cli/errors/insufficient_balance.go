package errors

import "bvpn-prototype/internal/infrastructure/errors"

func InsufficientBalanceError(data ...any) errors.Error {
	return errors.Error{
		Code: 11400,
		Type: "InsufficientBalance",
		Data: data,
	}
}
