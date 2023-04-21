package block_validators

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"bvpn-prototype/internal/protocol/interfaces"
	"bvpn-prototype/internal/protocol/protocol_error"
	"bvpn-prototype/internal/protocol/validators/node_validators"
	"bvpn-prototype/utils"
)

func offerValidation(block entity.Block, reader interfaces.ChainReader) error {
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

		// todo: what if somaone made offer not from himself
	}

	return nil
}
