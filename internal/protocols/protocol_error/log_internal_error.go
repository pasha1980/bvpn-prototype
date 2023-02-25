package protocol_error

const InternalErrorLog ErrorCode = 3

func LogInternalError(message string) *Error {
	return &Error{
		Code:    InternalErrorLog,
		Message: message,
	}
}
