package domain

import (
	"bvpn-prototype/internal/chain/api_out"
	chain_errors "bvpn-prototype/internal/chain/errors"
	"bvpn-prototype/internal/chain/storage"
	"bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/infrastructure/errors"
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
	GetChain(limit *int, offset *int) ([]entity.Block, error)
	ValidateStoredChain()
	UpdateChain()
}

type ChainServiceImpl struct {
	chainRepo   ChainRepository
	chainReader protocol.ChainReader
	mempoolRepo MempoolRepository
}

func (s *ChainServiceImpl) GetUTXOs() ([]block_data.ChainStored, error) {
	utxos, err := s.chainRepo.GetMyUTXOs()
	if err != nil {
		return nil, errors.StorageError(err.Error())
	}

	return utxos, nil
}

func (s *ChainServiceImpl) MakeNew(element block_data.ChainStored) (*block_data.ChainStored, error) {
	protocol.PrepareEntity(&element)
	s.mempoolRepo.AddNewElement(element)

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
			errors.StorageError(err.Error()).Log()
		}
	}
}

func (s *ChainServiceImpl) AddToMempool(element block_data.ChainStored, from *entity.Node) error {
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
	if !s.mempoolRepo.IsExist(element.ID) {
		s.mempoolRepo.AddNewElement(element)
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
		s.mempoolRepo.RemoveByIndex(datum.ID)
	}

	if block.Next == signer.GetAddr() {
		go s.createNewBlock()
	}

	return nil
}

func (s *ChainServiceImpl) GetChain(limit *int, offset *int) ([]entity.Block, error) {
	c, err := s.chainRepo.GetChain(limit, offset)
	if err != nil {
		return nil, errors.StorageError(err.Error())
	}

	return c, nil
}

func (s *ChainServiceImpl) ValidateStoredChain() {
	err := protocol.ValidateChain(s.chainReader)
	if err != nil {
		err.(errors.Error).Log()
	}
}

func (s *ChainServiceImpl) replaceChain(chain []entity.Block) error {
	err := s.validateChain(chain)
	if err != nil {
		return err
	}

	err = s.chainRepo.ReplaceChain(chain)
	if err != nil {
		return errors.StorageError(err.Error())
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
	if err != nil {
		return errors.StorageError(err.Error())
	}

	if lastBlock == nil {
		return chain_errors.EmptyChainError()
	}

	data := s.mempoolRepo.GetElements(protocol.BlockCapacity)

	newBlock, err := protocol.CreateNewBlock(*lastBlock, data)
	if err != nil {
		return err
	}

	if err = s.AddBlock(*newBlock, nil); err != nil {
		return errors.StorageError(err.Error())
	}

	for _, datum := range data {
		s.mempoolRepo.RemoveByIndex(datum.ID)
	}

	return nil
}

func NewChainService() (*ChainServiceImpl, error) {
	chainRepo, err := storage.NewChainRepo()
	if err != nil {
		return nil, errors.StorageError(err.Error())
	}

	mempoolRepo, err := storage.NewMempoolRepo()
	if err != nil {
		return nil, errors.StorageError(err.Error())
	}

	return &ChainServiceImpl{
		chainRepo:   chainRepo,
		chainReader: chainRepo,
		mempoolRepo: mempoolRepo,
	}, nil
}
