package protocol

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"bvpn-prototype/internal/protocol/hasher"
	"bvpn-prototype/internal/protocol/interfaces"
	"bvpn-prototype/internal/protocol/signer"
	"bvpn-prototype/internal/protocol/validators/block_validators"
	"github.com/google/uuid"
	"math"
	"math/rand"
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

	// todo: what if chain stopped

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
	rand.Seed(time.Now().Unix())

	block := reader.Last()
	if block == nil {
		return signer.GetAddr()
	}

	nodesTr := getLastTraffics(reader)

	var priorityList []string
	for addr, traffic := range nodesTr {
		for i := 0; i < traffic; i++ {
			priorityList = append(priorityList, addr)
		}
	}

	return priorityList[rand.Intn(len(priorityList))]
}

func getLastConnectionBreaks(reader interfaces.ChainReader) map[string]int {
	block := reader.Last()
	var connectionBreaks map[string]int
	breaksBorder := time.Now().Add(-1 * 30 * 24 * time.Hour)
	for {
		if block.TimeStamp.Before(breaksBorder) {
			break
		}

		for _, data := range block.Data {
			if data.Type != block_data.TypeConnectionBreak {
				continue
			}

			cb := data.Data.(*block_data.ConnectionBreak)
			savedCB, ok := connectionBreaks[cb.Node]
			if !ok {
				connectionBreaks[cb.Node] = 1
			} else {
				connectionBreaks[cb.Node] = savedCB + 1
			}
		}

		block = reader.Previous(block.Number)
	}

	return connectionBreaks
}

func getLastTraffics(reader interfaces.ChainReader) map[string]int {
	connectionBreaks := getLastConnectionBreaks(reader)

	block := reader.Last()
	var nodesTr map[string]int
	border := time.Now().Add(-1 * 10 * 24 * time.Hour)
	for {
		if block.TimeStamp.Before(border) {
			break
		}

		for _, data := range block.Data {
			if data.Type != block_data.TypeTraffic {
				continue
			}

			traffic := data.Data.(*block_data.Traffic)
			gb := int(math.Ceil(traffic.Bytes / 1073741824))
			savedTr, ok := nodesTr[traffic.Node]
			if !ok {
				nodesTr[traffic.Node] = gb
			} else {
				nodesTr[traffic.Node] = savedTr + gb
			}
		}

		block = reader.Previous(block.Number)
	}

	for addr, traffic := range nodesTr {
		degree, ok := connectionBreaks[addr]
		if ok {
			denominator := int(math.Pow(2, float64(degree)))
			nodesTr[addr] = traffic / denominator
		}

	}

	return nodesTr
}

func waitUntilMyTurn(previousBlockTime time.Time) {
	nextTimeStamp := previousBlockTime.Add(TimeToWait)
	time.Sleep(nextTimeStamp.Sub(time.Now()))
}
