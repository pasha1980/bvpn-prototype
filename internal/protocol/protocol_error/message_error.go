package protocol_error

const MessageErrorCode ErrorCode = 1

func MessageError(message string) *Error {
	return &Error{
		Code:    MessageErrorCode,
		Message: message,
	}
}
