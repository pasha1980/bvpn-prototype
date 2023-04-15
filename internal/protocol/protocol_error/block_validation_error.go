package protocol_error

import "strconv"

const BlockValidationErrorCode ErrorCode = 4

func BlockValidationError(message string, blockNumber uint64) *Error {
	return &Error{
		Code:    BlockValidationErrorCode,
		Message: message + "(block #" + strconv.FormatUint(blockNumber, 10) + ")",
	}
}
