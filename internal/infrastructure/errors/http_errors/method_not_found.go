package http_errors

import "bvpn-prototype/internal/infrastructure/errors"

func MethodNotFoundHttpError(data ...any) errors.Error {
	return errors.Error{
		Code: 10000,
		Type: "MethodNotFoundError",
		Data: data,
	}
}
