package mempool

import (
	"bvpn-prototype/internal/protocols/entity/block_data"
	"bvpn-prototype/internal/storage/config"
	utils2 "bvpn-prototype/utils"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"time"
)

func AddNewElement(element block_data.ChainStored) {
	updateData()
	defer storeData()
	storage.Data[element.ID.String()] = element
}

func GetElements(size int) []block_data.ChainStored {
	updateData()
	var result []block_data.ChainStored

	if utils2.SizeOf(storage.Data) < size {
		for _, value := range storage.Data {
			result = append(result, value)
		}
	}

	var usedIndexes []string
	var i int

	var temp []block_data.ChainStored
	for {
		element := randomElement(storage.Data)
		index := element.ID.String()
		if utils2.InStringSlice(index, usedIndexes) {
			continue
		}

		temp = append(temp, element)

		i = utils2.SizeOf(temp)
		if i > size {
			break
		}

		usedIndexes = append(usedIndexes, index)
		result = append(result, element)
	}

	return result
}

func IsExist(uuid uuid.UUID) bool {
	updateData()
	_, exist := storage.Data[uuid.String()]
	return exist
}

func RemoveByIndex(index uuid.UUID) {
	updateData()
	defer storeData()
	delete(storage.Data, index.String())
}

func randomElement(m map[string]block_data.ChainStored) block_data.ChainStored {
	var e block_data.ChainStored

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(m))
	for _, value := range m {
		if index == 0 {
			e = value
			break
		}
		index--
	}

	return e
}

func updateData() {
	storage = &pool{
		Data: make(map[string]block_data.ChainStored),
	}
	base64Encoded, err := os.ReadFile(config.Get().StorageDirectory + "/mempool.bvpn")
	if err != nil {
		storeData()
	}

	jsonEncoded, err := base64.StdEncoding.DecodeString(string(base64Encoded))
	if err != nil {
		storeData()
	}

	err = json.Unmarshal(jsonEncoded, &storage)
	if err != nil {
		storeData()
	}
}

func storeData() {
	jsonEncoded, _ := json.Marshal(storage)
	base64Encoded := base64.StdEncoding.EncodeToString(jsonEncoded)
	err := os.WriteFile(config.Get().StorageDirectory+"/mempool.bvpn", []byte(base64Encoded), 0666)
	if err != nil {
		// todo: what to do
	}
}
