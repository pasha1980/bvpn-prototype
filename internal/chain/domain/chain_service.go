package domain

import (
	"bvpn-prototype/internal/chain/api_out"
	"bvpn-prototype/internal/chain/storage"
	"bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/infrastructure/errors"
	"bvpn-prototype/internal/protocol"
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"bvpn-prototype/internal/protocol/interfaces"
	"bvpn-prototype/internal/protocol/params"
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
	chainReader interfaces.ChainReader
	mempoolRepo MempoolRepository
}

func (s *ChainServiceImpl) GetUTXOs() ([]block_data.ChainStored, error) {
	utxos, err := s.chainRepo.GetUTXOs(protocol.GetMyAddr())
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
	err := s.validateBlock(block)
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

func (s *ChainServiceImpl) GetMyLastOffer() (*block_data.Offer, error) {
	data, err := s.chainRepo.GetLastOffer(protocol.GetMyPubKey())
	if err != nil {
		return nil, errors.StorageError(err.Error())
	}

	if data == nil {
		return nil, nil
	}

	return data.Data.(*block_data.Offer), nil
}

func (s *ChainServiceImpl) SaveTraffic(traffic block_data.Traffic) error {
	element := block_data.ChainStored{
		Data: traffic,
		Type: block_data.TypeTraffic,
	}

	_, err := s.MakeNew(element)
	return err
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

func (s *ChainServiceImpl) validateBlock(block entity.Block) error {
	return protocol.ValidateBlock(block, s.chainReader)
}

func (s *ChainServiceImpl) createNewBlock() error {
	data := s.mempoolRepo.GetElements(params.BlockCapacity)

	newBlock, err := protocol.CreateNewBlock(s.chainReader, data)
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
