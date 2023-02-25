package block_validators

import (
	"bvpn-prototype/internal/protocols/entity"
	"errors"
	"time"
)

func timestampValidation(block entity.Block) error {
	if block.TimeStamp.After(time.Now()) {
		return errors.New("Invalid timestamp")
	}

	return nil
}
