package storage

import (
	"bvpn-prototype/internal/infrastructure/config"
	"bvpn-prototype/internal/protocol/entity/block_data"
	utils2 "bvpn-prototype/utils"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"time"
)

type MempoolRepository struct {
	Data map[string]block_data.ChainStored
}

func (r *MempoolRepository) AddNewElement(element block_data.ChainStored) {
	r.updateData()
	defer r.storeData()
	r.Data[element.ID.String()] = element
}

func (r *MempoolRepository) GetElements(size int) []block_data.ChainStored {
	r.updateData()
	var result []block_data.ChainStored

	if utils2.SizeOf(r.Data) < size {
		for _, value := range r.Data {
			result = append(result, value)
		}
	}

	var usedIndexes []string
	var i int

	var temp []block_data.ChainStored
	for {
		element := r.randomElement(r.Data)
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

func (r *MempoolRepository) IsExist(uuid uuid.UUID) bool {
	r.updateData()
	_, exist := r.Data[uuid.String()]
	return exist
}

func (r *MempoolRepository) RemoveByIndex(index uuid.UUID) {
	r.updateData()
	defer r.storeData()
	delete(r.Data, index.String())
}

func (r *MempoolRepository) randomElement(m map[string]block_data.ChainStored) block_data.ChainStored {
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

func (r *MempoolRepository) updateData() {
	r.Data = make(map[string]block_data.ChainStored)
	base64Encoded, err := os.ReadFile(config.Get().StorageDirectory + "/mempool.bvpn")
	if err != nil {
		r.storeData() // todo
	}

	jsonEncoded, err := base64.StdEncoding.DecodeString(string(base64Encoded))
	if err != nil {
		r.storeData()
	}

	err = json.Unmarshal(jsonEncoded, &r.Data)
	if err != nil {
		r.storeData()
	}
}

func (r *MempoolRepository) storeData() {
	jsonEncoded, _ := json.Marshal(r.Data)
	base64Encoded := base64.StdEncoding.EncodeToString(jsonEncoded)
	err := os.WriteFile(config.Get().StorageDirectory+"/mempool.bvpn", []byte(base64Encoded), 0666)
	if err != nil {
		// todo: what to do
	}
}

func NewMempoolRepo() (*MempoolRepository, error) {
	repo := MempoolRepository{}
	repo.updateData()
	return &repo, nil
}
