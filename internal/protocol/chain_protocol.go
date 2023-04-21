package protocol

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"bvpn-prototype/internal/protocol/hasher"
	"bvpn-prototype/internal/protocol/signer"
	"bvpn-prototype/internal/protocol/validators/block_validators"
	"github.com/google/uuid"
	"time"
)

const BlockCapacity = 1048576
const TimeToWait = 10 * time.Second

func ValidateChain(reader ChainReader) error {
	block := reader.Last()
	if block == nil {
		// todo
	}

	reader.Start()
	previousBlock := reader.Next()

	for {
		block = reader.Next()
		if block == nil {
			break
		}

		err := ValidateBlock(*block, previousBlock)
		if err != nil {
			return err
		}

		previousBlock = block
	}

	return nil
}

// todo
func AddInitialBlock() error {
	timestamp, _ := time.Parse("2006-01-02 15:04:05", entity.InitialBlockTimestamp)

	initialBlock := entity.Block{
		Number:       1,
		PreviousHash: entity.InitialBlockPrevHash,
		TimeStamp:    timestamp,
		Data:         []block_data.ChainStored{},
		CreatedBy:    "0",
		Next:         "0",
	}

	initialBlock.Hash = hasher.EncryptBlock(initialBlock)

	//err := p.AddBlock(initialBlock, nil)
	//if err != nil {
	//	return err
	//}

	return nil
}

func ValidateBlock(block entity.Block, previousBlock *entity.Block) error {
	var err error

	for _, validator := range block_validators.GetValidationRules() {
		err = validator(block, previousBlock)
	}

	return err
}

func CreateNewBlock(previousBlock entity.Block, data []block_data.ChainStored) (*entity.Block, error) {

	waitUntilMyTurn(previousBlock.TimeStamp)

	newBlock := entity.Block{
		Number:       previousBlock.Number + 1,
		PreviousHash: previousBlock.PreviousHash,
		Data:         data,
		TimeStamp:    time.Now(),
		CreatedBy:    signer.GetAddr(),
		Next:         whoIsNext(),
	}
	newBlock.Hash = hasher.EncryptBlock(newBlock)

	err := ValidateBlock(newBlock, &previousBlock)
	if err != nil {
		return nil, err
	}

	return &newBlock, nil
}

func PrepareEntity(data *block_data.ChainStored) {
	data.ID = uuid.New()
	signer.Sign(data)
}

func InitKeys() {
	signer.Init()
}

func GetMyAddr() string {
	return signer.GetAddr()
}

func GetMyPubKey() string {
	return signer.GetPubKey()
}

func whoIsNext() string {
	// todo
	return ""
}

func waitUntilMyTurn(previousBlockTime time.Time) {
	nextTimeStamp := previousBlockTime.Add(TimeToWait)
	time.Sleep(nextTimeStamp.Sub(time.Now()))
}
