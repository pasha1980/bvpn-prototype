package protocol_error

const LogErrorCode ErrorCode = 2

func LogError(message string) *Error {
	return &Error{
		Code:    LogErrorCode,
		Message: message,
	}
}
