package protocol_error

import "bvpn-prototype/internal/infrastructure/errors"

func PeerValidationError(message string) errors.Error {
	return errors.Error{
		Code: 12001,
		Type: "PeerValidationError",
		Data: message,
	}
}
