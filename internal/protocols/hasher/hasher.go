package hasher

import (
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/entity/block_data"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"
)

/*
 Hash Algorithm: SHA-256
*/

type BlockToEncrypt struct {
	Number       uint64                   `json:"number"`
	PreviousHash string                   `json:"previous_hash"`
	Data         []block_data.ChainStored `json:"data"`
	TimeStamp    time.Time                `json:"time_stamp"`
}

func EncryptBlock(block entity.Block) []byte {
	numStr := strconv.FormatUint(block.Number, 10)

	dataJson, _ := json.Marshal(block.Data)
	dataStr := base64.StdEncoding.EncodeToString(dataJson)

	timeStr := strconv.Itoa(int(block.TimeStamp.Unix()))

	blockStr := numStr + block.PreviousHash + dataStr + timeStr

	return encryptRaw(blockStr)
}

func encryptRaw(data string) []byte {
	h := sha256.New()
	h.Write([]byte(data))
	return h.Sum(nil)
}
