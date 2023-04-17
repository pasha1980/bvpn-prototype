package protocol_error

import (
	"bvpn-prototype/internal/infrastructure/errors"
)

func BlockValidationError(message string, blockNumber uint64) errors.Error {
	return errors.Error{
		Code: 13001,
		Type: "BlockValidationError",
		Data: map[string]any{
			"message": message,
			"block":   blockNumber,
		},
	}
}
