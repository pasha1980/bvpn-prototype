package errors

import "bvpn-prototype/internal/infrastructure/errors"

func NoOfferError(data ...any) errors.Error {
	return errors.Error{
		Code: 11300,
		Type: "NoOffer",
		Data: data,
	}
}
