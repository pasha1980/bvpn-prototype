package common

import "bvpn-prototype/internal/infrastructure/errors"

func FilesystemError(data ...any) errors.Error {
	return errors.Error{
		Code: 10011,
		Type: "FilesystemError",
		Data: data,
	}
}
