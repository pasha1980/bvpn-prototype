package protocols

import (
	"bvpn-prototype/internal/http/http_out"
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/entity/block_data"
	"bvpn-prototype/internal/protocols/hasher"
	"bvpn-prototype/internal/protocols/protocol_error"
	"bvpn-prototype/internal/protocols/repo"
	"bvpn-prototype/internal/protocols/signer"
	"bvpn-prototype/internal/protocols/validators/block_validators"
	"bvpn-prototype/internal/storage/chain"
	"bvpn-prototype/internal/storage/mempool"
	"github.com/google/uuid"
	"sort"
	"time"
)

type ChainProtocol struct {
	repo repo.ChainStorageRepo
}

func (p *ChainProtocol) New(data block_data.ChainStored) block_data.ChainStored {
	data.ID = uuid.New()
	signer.Sign(&data)
	mempool.AddNewElement(data)
	http_out.BroadcastMempool(data, GetPeerProtocol().GetPeers())
	return data
}

func (p *ChainProtocol) GetUTXOs() ([]block_data.ChainStored, error) {
	utxos, err := p.repo.GetMyUTXOs()
	if err != nil {
		return nil, protocol_error.LogInternalError(err.Error())
	}

	return utxos, nil
}

func (p *ChainProtocol) ValidateChain() error {
	block, err := p.repo.GetLastBlock()
	if err != nil {
		return protocol_error.LogInternalError(err.Error())
	}

	if block == nil {
		return p.AddInitialBlock()
	}

	for {
		if block.PreviousHash == entity.InitialBlockPrevHash {
			break
		}

		previousBlock, err := p.repo.GetBlockByHash(block.PreviousHash)
		if err != nil {
			return protocol_error.LogInternalError(err.Error())
		}

		err = p.ValidateBlock(*block, previousBlock)
		if err != nil {
			return err
		}

		block = previousBlock
	}

	return nil
}

func (p *ChainProtocol) UpdateChain() {
	var err error
	var chains [][]entity.Block

	for _, node := range GetPeerProtocol().GetPeers() {
		peerChain := http_out.GetFullChain(node)
		chains = append(chains, peerChain)
	}

	if len(chains) == 0 {
		return
	}

	for _, peerChain := range chains {
		err = p.validateGivenChain(peerChain)
		if err != nil {
			continue
		}

		err = p.ReplaceChain(peerChain)
		if err != nil {
			continue
		}
	}
}

func (p *ChainProtocol) ReplaceChain(chain []entity.Block) error {
	err := p.validateGivenChain(chain)
	if err != nil {
		return err
	}

	err = p.repo.ReplaceChain(chain)
	if err != nil {
		return protocol_error.LogInternalError(err.Error())
	}

	return nil
}

func (p *ChainProtocol) AddBlock(block entity.Block) error {
	lastBlock, err := p.repo.GetLastBlock()
	if err != nil {
		return protocol_error.LogInternalError(err.Error())
	}

	err = p.ValidateBlock(block, lastBlock)
	if err != nil {
		return err
	}

	http_out.BroadcastBlock(block, GetPeerProtocol().GetPeers())

	_, err = p.repo.SaveBlock(block)
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

	initialBlock.Hash = hasher.EncryptBlock(initialBlock)
	err := p.AddBlock(initialBlock)
	if err != nil {
		return err
	}

	return nil
}

func (p *ChainProtocol) ValidateBlock(block entity.Block, previousBlock *entity.Block) error {
	var err error

	for _, validator := range block_validators.GetValidationRules() {
		err = validator(block, previousBlock)
	}

	return err
}

func (p *ChainProtocol) AddToMempool(element block_data.ChainStored) {
	if !mempool.IsExist(element.ID) {
		mempool.AddNewElement(element)

		http_out.BroadcastMempool(element, GetPeerProtocol().GetPeers())
	}
}

func (p *ChainProtocol) GetChain(limit int, offset int) ([]entity.Block, error) {
	c, err := p.repo.GetChain(limit, offset)
	if err != nil {
		return nil, protocol_error.LogInternalError(err.Error())
	}

	return c, nil
}

func (p *ChainProtocol) validateGivenChain(chain []entity.Block) error {
	if len(chain) == 0 {
		return protocol_error.MessageError("Empty chain")
	}

	sort.Slice(chain, func(i, j int) bool {
		return chain[i].Number < chain[j].Number
	})

	for i, block := range chain {
		if block.IsInitial() {
			break
		}

		err := p.ValidateBlock(block, &chain[i+1])
		if err != nil {
			return err
		}

	}

	lastBlock, err := p.repo.GetLastBlock()
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

func GetChainProtocol() *ChainProtocol {
	return &ChainProtocol{
		repo: chain.NewChainRepo(),
	}
}
