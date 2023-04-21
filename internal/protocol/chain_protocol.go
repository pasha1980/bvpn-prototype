package protocol

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"bvpn-prototype/internal/protocol/hasher"
	"bvpn-prototype/internal/protocol/interfaces"
	"bvpn-prototype/internal/protocol/signer"
	"bvpn-prototype/internal/protocol/validators/block_validators"
	"github.com/google/uuid"
	"time"
)

const BlockCapacity = 1048576
const TimeToWait = 10 * time.Second

func ValidateChain(reader interfaces.ChainReader) error {
	block := reader.Last()
	if block == nil {
		// todo
	}

	reader.Start()

	for {
		block = reader.Next()
		if block == nil {
			break
		}

		err := ValidateBlock(*block, reader)
		if err != nil {
			return err
		}
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

func ValidateBlock(block entity.Block, chainReader interfaces.ChainReader) error {
	var err error

	for _, validator := range block_validators.GetValidationRules() {
		err = validator(block, chainReader)
	}

	return err
}

func CreateNewBlock(chainReader interfaces.ChainReader, data []block_data.ChainStored) (*entity.Block, error) {
	lastBlock := chainReader.Last()
	waitUntilMyTurn(lastBlock.TimeStamp)

	newBlock := entity.Block{
		Number:       lastBlock.Number + 1,
		PreviousHash: lastBlock.PreviousHash,
		Data:         data,
		TimeStamp:    time.Now(),
		CreatedBy:    signer.GetAddr(),
		Next:         whoIsNext(chainReader),
	}
	newBlock.Hash = hasher.EncryptBlock(newBlock)

	err := ValidateBlock(newBlock, chainReader)
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

func whoIsNext(reader interfaces.ChainReader) string {
	// todo
	return ""
}

func waitUntilMyTurn(previousBlockTime time.Time) {
	nextTimeStamp := previousBlockTime.Add(TimeToWait)
	time.Sleep(nextTimeStamp.Sub(time.Now()))
}
