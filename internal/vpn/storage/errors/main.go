package errors

type profileStorageError string

func (e profileStorageError) Error() string {
	return string(e)
}
