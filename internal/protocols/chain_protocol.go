package protocols

import (
	"bvpn-prototype/internal/http/http_out"
	"bvpn-prototype/internal/mempool"
	"bvpn-prototype/internal/protocols/block_validators"
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/entity/block_data"
	"bvpn-prototype/internal/protocols/hasher"
	"bvpn-prototype/internal/protocols/protocol_error"
	"bvpn-prototype/internal/protocols/repo"
	"errors"
	"sort"
	"strconv"
	"time"
)

/*

Chain methods:
+ ValidateChain() error
+ UpdateChain() error
+ ReplaceChain(chain []Block) error

BlockMethods:
- CreateNewBlock()
+ AddBlock(block Block)
+ ValidateBlock(block Block)
+ ValidateBlockData(data) error

MempoolMethods:
+ AddToMempool(data)
*/

type ChainProtocol struct {
	nodes     []entity.Node
	chainRepo repo.ChainStorageRepo
}

func (p *ChainProtocol) ValidateChain() error {
	block, err := p.chainRepo.GetLastBlock()
	if err != nil {
		// todo: log error
	}

	if block == nil {
		return p.AddInitialBlock()
	}

	for {
		err := p.ValidateBlock(*block)
		if err != nil {
			return err
		}

		if block.PreviousHash == entity.InitialBlockPrevHash {
			break
		}

		block, err = p.chainRepo.GetBlockByHash(block.Hash)
		if err != nil {
			return protocol_error.MessageError("Storage error")
		}

		if block == nil {
			return protocol_error.MessageError("Invalid previous hash on block #" + strconv.FormatUint(block.Number, 10))
		}
	}

	return nil
}

func (p *ChainProtocol) UpdateChain() error {
	var err error
	var chains [][]entity.Block

	for _, node := range p.nodes {
		chain := http_out.GetFullChain(node)
		chains = append(chains, chain)
	}

	if len(chains) == 0 {
		return nil
	}

	for _, chain := range chains {
		err = p.validateGivenChain(chain)
		if err != nil {
			continue
		}

		err = p.ReplaceChain(chain)
		if err != nil {
			continue
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
		return protocol_error.LogInternalError(err.Error())
	}

	return nil
}

func (p *ChainProtocol) AddBlock(block entity.Block) error {
	err := p.ValidateBlock(block)
	if err != nil {
		return err
	}

	lastBlock, err := p.chainRepo.GetLastBlock()
	if err != nil {
		return protocol_error.LogInternalError(err.Error())
	}

	if lastBlock == nil {
		return p.AddInitialBlock()
	}

	if block.Number != lastBlock.Number+1 || block.PreviousHash != lastBlock.Hash {
		return protocol_error.MessageError("Invalid block")
	}

	http_out.BroadcastBlock(block, p.nodes)

	_, err = p.chainRepo.SaveBlock(block)
	if err != nil {
		return protocol_error.LogInternalError(err.Error())
	}

	return nil
}

func (p *ChainProtocol) AddInitialBlock() error {
	timestamp, _ := time.Parse("2006-01-02 15:04:05", entity.InitialBlockTimestamp)

	initialBlock := entity.Block{
		Number:       1,
		PreviousHash: entity.InitialBlockPrevHash,
		TimeStamp:    timestamp,
		Data:         []block_data.ChainStored{},
	}

	initialBlock.Hash = string(hasher.EncryptBlock(initialBlock))
	_, err := p.chainRepo.SaveBlock(initialBlock)
	if err != nil {
		return protocol_error.LogInternalError(err.Error())
	}

	return nil
}

func (p *ChainProtocol) ValidateBlock(block entity.Block) error {
	var err error

	for _, validator := range block_validators.GetValidationRules() {
		err = validator(block)
	}

	return err
}

func (p *ChainProtocol) AddToMempool(element block_data.ChainStored) {
	if !mempool.IsExist(element.ID) {
		mempool.AddNewElement(element)

		http_out.BroadcastMempool(element, p.nodes)
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

	lastBlock, err := p.chainRepo.GetLastBlock()
	if err != nil {
		return protocol_error.LogInternalError(err.Error())
	}

	if lastBlock != nil {
		if len(chain) <= int(lastBlock.Number) {
			return protocol_error.MessageError("Chain is not prioritized")
		}
	}

	return nil
}
