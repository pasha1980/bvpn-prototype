package protocol_error

const ErrorLog ErrorCode = 2

func LogError(message string) *Error {
	return &Error{
		Code:    ErrorLog,
		Message: message,
	}
}
