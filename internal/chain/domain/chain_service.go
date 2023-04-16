package domain

import (
	"bvpn-prototype/internal/chain/api_out"
	"bvpn-prototype/internal/chain/mempool"
	"bvpn-prototype/internal/chain/storage"
	"bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/protocol"
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"bvpn-prototype/internal/protocol/signer"
	"bvpn-prototype/internal/protocol/validators/node_validators"
	"bvpn-prototype/utils"
)

type ChainService interface {
	AddToMempool(entity block_data.ChainStored, from *entity.Node) error
	AddBlock(block entity.Block, from *entity.Node) error
	GetChain(limit int, offset int) ([]entity.Block, error)
}

type ChainServiceImpl struct {
	chainRepo   ChainRepository
	chainReader protocol.ChainReader
}

func (s *ChainServiceImpl) GetUTXOs() ([]block_data.ChainStored, error) {
	utxos, err := s.chainRepo.GetMyUTXOs()
	if err != nil {
		return nil, err // todo: domain errors
	}

	return utxos, nil
}

func (*ChainServiceImpl) MakeNew(element block_data.ChainStored) (*block_data.ChainStored, error) {
	protocol.PrepareEntity(&element)
	mempool.AddNewElement(element) // todo

	peers := di.Get("peer_public").(PeerPublicService).GetPeers(nil)
	go api_out.BroadcastMempool(element, peers) // todo
	return &element, nil
}

func (s *ChainServiceImpl) UpdateChain() {
	var err error
	var chains [][]entity.Block

	peers := di.Get("peer_public").(PeerPublicService).GetPeers(nil)
	for _, node := range peers {
		peerChain := api_out.GetFullChain(node) // todo
		chains = append(chains, peerChain)
	}

	if len(chains) == 0 {
		return
	}

	for _, peerChain := range chains {
		err = s.validateChain(peerChain)
		if err != nil {
			continue
		}

		err = s.replaceChain(peerChain)
		if err != nil {
			continue
		}
	}
}

func (*ChainServiceImpl) AddToMempool(element block_data.ChainStored, from *entity.Node) error {
	if element.Type == block_data.TypeOffer {
		url := element.Data.(block_data.Offer).URL
		node := entity.Node{
			URL: url,
			IP:  utils.GetIp(url),
		}

		rules := node_validators.GetValidationRules()
		for _, rule := range rules {
			err := rule(node)
			if err != nil {
				return err
			}
		}
	}

	peers := di.Get("peer_public").(PeerPublicService).GetPeers(from)
	// todo
	if !mempool.IsExist(element.ID) {
		mempool.AddNewElement(element)
		go api_out.BroadcastMempool(element, peers)
	}

	return nil
}

func (s *ChainServiceImpl) AddBlock(block entity.Block, from *entity.Node) error {
	lastBlock, err := s.chainRepo.GetLastBlock()
	if err != nil {
		return err // todo: domain error
	}

	err = s.validateBlock(block, lastBlock) // todo: what if initial
	if err != nil {
		return err
	}

	peers := di.Get("peer_public").(PeerPublicService).GetPeers(from)
	go api_out.BroadcastBlock(block, peers)

	_, err = s.chainRepo.SaveBlock(block)
	if err != nil {
		return err // todo: domain error
	}

	for _, datum := range block.Data {
		mempool.RemoveByIndex(datum.ID)
	}

	if block.Next == signer.GetAddr() {
		go s.createNewBlock()
	}

	return nil
}

func (s *ChainServiceImpl) GetChain(limit int, offset int) ([]entity.Block, error) {
	c, err := s.chainRepo.GetChain(limit, offset)
	if err != nil {
		return nil, err // todo: domain errors
	}

	return c, nil
}

func (s *ChainServiceImpl) ValidateStoredChain() error {
	return protocol.ValidateChain(s.chainReader)
}

func (s *ChainServiceImpl) replaceChain(chain []entity.Block) error {
	err := s.validateChain(chain)
	if err != nil {
		return err
	}

	err = s.chainRepo.ReplaceChain(chain)
	if err != nil {
		return err // todo: domain errors
	}

	return nil
}

func (*ChainServiceImpl) validateChain(chain []entity.Block) error {
	return protocol.ValidateChain(NewSliceChainReader(chain))
}

func (*ChainServiceImpl) validateBlock(block entity.Block, previousBlock *entity.Block) error {
	return protocol.ValidateBlock(block, previousBlock)
}

func (s *ChainServiceImpl) createNewBlock() error {
	lastBlock, err := s.chainRepo.GetLastBlock()
	if err != nil || lastBlock == nil {
		return err // todo: domain errors
	}

	data := mempool.GetElements(protocol.BlockCapacity)

	newBlock, err := protocol.CreateNewBlock(*lastBlock, data)
	if err != nil {
		return err
	}

	if err = s.AddBlock(*newBlock, nil); err != nil {
		return err
	}

	for _, datum := range data {
		mempool.RemoveByIndex(datum.ID)
	}

	return nil
}

func NewChainService() (*ChainServiceImpl, error) {
	chainRepo, err := storage.NewChainRepo()
	if err != nil {
		return nil, err
	}

	return &ChainServiceImpl{
		chainRepo:   chainRepo,
		chainReader: chainRepo,
	}, nil
}
