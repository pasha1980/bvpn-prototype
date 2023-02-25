package protocol_error

const PeerValidationErrorCode ErrorCode = 5

func PeerValidationError(message string) *Error {
	return &Error{
		Code:    PeerValidationErrorCode,
		Message: message,
	}
}
