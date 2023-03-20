package block_validators

import (
	"bvpn-prototype/internal/http/http_out"
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/entity/block_data"
	"bvpn-prototype/internal/protocols/hasher"
	"bvpn-prototype/internal/protocols/protocol_error"
	"bvpn-prototype/internal/protocols/validators/node_validators"
	"bvpn-prototype/internal/utils"
)

func offerValidation(block entity.Block, previousBlock *entity.Block) error {
	for _, data := range block.Data {
		if data.Type != block_data.TypeOffer {
			continue
		}

		url := data.Data.(block_data.Offer).URL
		node := entity.Node{
			URL: url,
			IP:  utils.GetIp(url),
		}

		rules := node_validators.GetValidationRules()
		for _, rule := range rules {
			err := rule(node)
			if err != nil {
				return protocol_error.BlockValidationError("Invalid offer ("+err.Error()+")", block.Number)
			}
		}

		address, err := http_out.GetAddr(node)
		if err != nil {
			return protocol_error.BlockValidationError("Invalid offer ("+err.Error()+")", block.Number)
		}

		if address != hasher.EncryptString(data.PubKey) {
			return protocol_error.BlockValidationError("Invalid offer (invalid signature)", block.Number)
		}
	}

	return nil
}
