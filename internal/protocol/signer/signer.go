package signer

import (
	"bvpn-prototype/internal/logger"
	"bvpn-prototype/internal/protocols/entity/block_data"
	"bvpn-prototype/internal/protocols/hasher"
	"bvpn-prototype/internal/storage/config"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/sha3"
	"os"
)

type s struct {
	pub ecdsa.PublicKey
	prv ecdsa.PrivateKey
}

func (st *s) save() {
	dir := config.Get().StorageDirectory
	file, err := os.Create(dir + "/prv.pem")
	if err != nil {
		logger.LogError(err.Error())
	}

	encoded, _ := x509.MarshalECPrivateKey(&st.prv)
	pem.Encode(file, &pem.Block{
		Type:    "BVPN PRIVATE KEY",
		Headers: nil,
		Bytes:   encoded,
	})
}

func storage() *s {
	var st s
	dir := config.Get().StorageDirectory
	data, err := os.ReadFile(dir + "/prv.pem")
	if err != nil {
		logger.LogError(err.Error())
	}
	block, _ := pem.Decode(data)
	encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(encoded)
	st.prv = *privateKey
	st.pub = st.prv.PublicKey
	return &st
}

func Init() {
	var st s
	prv, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	st.prv = *prv
	st.save()
}

func Validate(data *block_data.ChainStored) bool {
	decodedPub, _ := x509.ParsePKIXPublicKey([]byte(data.PubKey))
	pub := decodedPub.(*ecdsa.PublicKey)
	hash := []byte(hasher.EncryptString(fmt.Sprintf("%v", *data)))
	return ecdsa.VerifyASN1(pub, hash[:], []byte(data.Sign))
}

func Sign(data *block_data.ChainStored) {
	st := storage()
	hash := []byte(hasher.EncryptString(fmt.Sprintf("%v", *data)))
	sign, _ := ecdsa.SignASN1(rand.Reader, &st.prv, hash)
	data.Sign = fmt.Sprintf("%x", sign)
	encodedPub, _ := x509.MarshalPKIXPublicKey(&st.pub)
	data.PubKey = fmt.Sprintf("%x", encodedPub)
}

func GetAddr() string {
	st := storage()
	buf, _ := x509.MarshalPKIXPublicKey(&st.pub)
	h := make([]byte, 32)
	sha3.ShakeSum128(h, buf)
	return fmt.Sprintf("%x", h)
}
