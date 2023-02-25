package protocol_error

import "bvpn-prototype/internal/logger"

func Handle(err *Error) (string, int) {
	switch err.Code {
	case MessageErrorCode, PeerValidationErrorCode, BlockValidationErrorCode:
		return err.Message, 400
	case LogErrorCode:
		logger.LogError(err.Error())
		return err.Message, 400
	case LogInternalErrorCode:
		logger.LogError(err.Error())
		return err.Message, 500
	}

	return "UndefinedError", 200
}
