package protocol_error

import "strconv"

type Error struct {
	Code    ErrorCode
	Message string
}

type ErrorCode int

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return strconv.Itoa(int(e.Code))
}
