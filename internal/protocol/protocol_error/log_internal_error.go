package protocol_error

const LogInternalErrorCode ErrorCode = 3

func LogInternalError(message string) *Error {
	return &Error{
		Code:    LogInternalErrorCode,
		Message: message,
	}
}
