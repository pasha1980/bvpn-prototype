package errors

import "bvpn-prototype/internal/infrastructure/errors"

func PeerNotAvailable(data ...any) errors.Error {
	return errors.Error{
		Code: 11201,
		Type: "PeerNotAvailable",
		Data: data,
	}
}
