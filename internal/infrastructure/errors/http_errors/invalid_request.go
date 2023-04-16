package http_errors

import "bvpn-prototype/internal/infrastructure/errors"

func InvalidRequest(data ...any) errors.Error {
	return errors.Error{
		Code: 10000,
		Type: "InvalidRequest",
		Data: data,
	}
}
