package hasher

import (
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/entity/block_data"
	"fmt"
	"golang.org/x/crypto/sha3"
	"strconv"
	"time"
)

/*
 Hash Algorithm: SHA-256
*/

type BlockToEncrypt struct {
	Number       uint64
	PreviousHash string
	Data         []block_data.ChainStored
	TimeStamp    time.Time
}

func EncryptBlock(block entity.Block) []byte {
	numStr := strconv.FormatUint(block.Number, 10)
	dataStr := fmt.Sprintf("%v", block.Data)
	timeStr := strconv.Itoa(int(block.TimeStamp.Unix()))
	blockStr := numStr + block.PreviousHash + dataStr + timeStr
	return EncryptString(blockStr)
}

func EncryptString(data string) []byte {
	hash := sha3.New256()
	hash.Write([]byte(data))
	return hash.Sum(nil)
}
