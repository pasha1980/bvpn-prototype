package errors

func StorageError(data ...any) Error {
	return Error{
		Code: 10011,
		Type: "StorageError",
		Data: data,
	}
}
