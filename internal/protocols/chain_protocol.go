package protocols

import (
	"bvpn-prototype/internal/mempool"
	"bvpn-prototype/internal/protocols/data_validators"
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/entity/block_data"
	"bvpn-prototype/internal/protocols/repo"
	"errors"
	"sort"
	"strconv"
)

/*

Chain methods:
+ ValidateChain() error
+ UpdateChain() error
+ ReplaceChain(chain []Block) error

BlockMethods:
- CreateNewBlock()
- AddBlock(block Block)
+ ValidateBlock(block Block)
- ValidateBlockData(data) error

MempoolMethods:
+ AddToMempool(data)
*/

type ChainProtocol struct {
	chainRepo repo.ChainStorageRepo
}

func (p *ChainProtocol) ValidateChain() error {
	block := p.chainRepo.GetLastBlock()
	if block == nil {
		return nil
	}

	for {
		err := p.ValidateBlock(*block)
		if err != nil {
			return err
		}

		if block.PreviousHash == "" {
			break
		}

		block = p.chainRepo.GetBlockByHash(block.Hash)
		if block == nil {
			return errors.New("Invalid previous hash on block #" + strconv.FormatUint(block.Number, 10)) // todo: Custom error
		}
	}

	return nil
}

func (p *ChainProtocol) UpdateChain() error {
	var err error
	var chains [][]entity.Block

	// todo: Get chains from nodes

	for _, chain := range chains {
		err = p.validateGivenChain(chain)
		if err != nil {
			continue
		}

		err = p.ReplaceChain(chain)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ChainProtocol) ReplaceChain(chain []entity.Block) error {
	err := p.validateGivenChain(chain)
	if err != nil {
		return err
	}

	err = p.chainRepo.ReplaceChain(chain)
	if err != nil {
		return err
	}

	return nil
}

func (p *ChainProtocol) ValidateBlock(block entity.Block) error {
	var err error

	for _, validator := range data_validators.GetValidationRules() {
		err = validator(block)
	}

	return err
}

func (p *ChainProtocol) AddToMempool(element block_data.ChainStored) {
	if !mempool.IsExist(element.ID) {
		mempool.AddNewElement(element)

		// todo: send to next node
	}
}

func (p *ChainProtocol) validateGivenChain(chain []entity.Block) error {
	if len(chain) == 0 {
		return errors.New("Empty chain") // todo: Custom error
	}

	sort.Slice(chain, func(i, j int) bool {
		return chain[i].Number < chain[j].Number
	})

	for i, block := range chain {
		err := p.ValidateBlock(block)
		if err != nil {
			return err
		}

		if block.IsInitial() {
			break
		}

		nextBlock := chain[i+1]
		if block.PreviousHash != nextBlock.Hash || block.Number != nextBlock.Number+1 {
			return errors.New("Disconnected chain") // todo: Custom error
		}
	}

	lastBlock := p.chainRepo.GetLastBlock()
	if lastBlock != nil {
		if len(chain) <= int(lastBlock.Number) {
			return errors.New("Chain is not prioritized") // todo: Custom error
		}
	}

	return nil
}
