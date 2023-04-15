package http_errors

const (
	ErrorMethodNotFound  HttpErrors = "MethodNotFound"
	ErrorNotEnoughParams HttpErrors = "NotEnoughParams"
	ErrorInvalidParams   HttpErrors = "InvalidParams"
	ErrorInvalidRequest  HttpErrors = "InvalidRequest"
)

type HttpErrors string

func (e HttpErrors) Error() string {
	return string(e)
}
