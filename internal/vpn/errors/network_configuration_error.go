package errors

import "bvpn-prototype/internal/infrastructure/errors"

func NetworkConfigurationError(data ...any) errors.Error {
	return errors.Error{
		Code: 11311,
		Type: "NetworkConfigurationError",
		Data: data,
	}
}
