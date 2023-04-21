package storage

import (
	"encoding/base64"
	"encoding/json"
	"io"
)

type ClientStorageFormat struct {
	Id     string `json:"id"`
	PrvKey string `json:"prvKey"`
	PubKey string `json:"pubKey"`
	Client string `json:"client"`

	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
	URL       string  `json:"URL"`
}

func Read(reader io.Reader) (*ClientStorageFormat, error) {
	var format ClientStorageFormat

	var base64Encoded []byte
	_, err := reader.Read(base64Encoded)
	if err != nil {
		return nil, err
	}

	jsonEncoded, err := base64.StdEncoding.DecodeString(string(base64Encoded))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonEncoded, &format)
	if err != nil {
		return nil, err
	}

	return &format, nil
}

func (k *ClientStorageFormat) Write(writer io.Writer) error {
	jsonEncoded, _ := json.Marshal(k)
	var encoded []byte
	base64.StdEncoding.Encode(encoded, jsonEncoded)
	_, err := writer.Write(encoded)
	if err != nil {
		return err
	}
	return err
}
