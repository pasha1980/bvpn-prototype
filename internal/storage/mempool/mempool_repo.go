package mempool

import (
	"bvpn-prototype/internal/protocols/entity/block_data"
	"bvpn-prototype/internal/utils"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

func AddNewElement(element block_data.ChainStored) {
	storage[element.ID.String()] = element
}

func GetRandomElements(count int) []block_data.ChainStored {
	var result []block_data.ChainStored

	if len(storage) < count {
		for _, value := range storage {
			result = append(result, value)
		}
	}

	var usedIndexes []string
	var i int

	for {
		element := randomElement(storage)
		index := element.ID.String()
		if utils.InStringSlice(index, usedIndexes) {
			continue
		}

		usedIndexes = append(usedIndexes, index)
		result = append(result, element)

		i++
		if i == count {
			break
		}
	}

	return result
}

func IsExist(uuid uuid.UUID) bool {
	_, exist := storage[uuid.String()]
	return exist
}

func RemoveByIndex(index uuid.UUID) {
	delete(storage, index.String())
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
