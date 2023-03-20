package hasher

import (
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/entity/block_data"
	"fmt"
	"golang.org/x/crypto/sha3"
	"sort"
	"strconv"
	"time"
)

/*
 Hash Algorithm: SHA-3-256
*/

type BlockToEncrypt struct {
	Number       uint64
	PreviousHash string
	Data         []block_data.ChainStored
	TimeStamp    time.Time
}

func EncryptBlock(block entity.Block) string {
	numStr := strconv.FormatUint(block.Number, 10)

	sort.Slice(block.Data, func(i, j int) bool {
		return block.Data[i].ID.ID() < block.Data[j].ID.ID()
	})
	dataStr := fmt.Sprintf("%v", block.Data)

	timeStr := strconv.Itoa(int(block.TimeStamp.Unix()))
	blockStr := numStr + block.PreviousHash + dataStr + timeStr
	return EncryptString(blockStr)
}

func EncryptString(data string) string {
	s := sha3.New256()
	s.Write([]byte(data))
	h := s.Sum(nil)
	return fmt.Sprintf("%x", h)
}
