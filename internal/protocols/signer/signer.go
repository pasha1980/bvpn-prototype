package signer

import (
	"bvpn-prototype/internal/logger"
	"bvpn-prototype/internal/protocols/entity/block_data"
	"bvpn-prototype/internal/protocols/hasher"
	"bvpn-prototype/internal/storage/config"
	"fmt"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/curve"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/ecdsa"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/privatekey"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/publickey"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/signature"
	"math/big"
	"math/rand"
	"os"
	"time"
)

type s struct {
	addr string
	pub  publickey.PublicKey
	prv  privatekey.PrivateKey
}

func (st *s) save() {
	dir := config.Get().StorageDirectory
	err := os.WriteFile(dir+"/prv.pem", []byte(st.prv.ToPem()), 0600)
	if err != nil {
		logger.LogError(err.Error())
	}
}

func storage() *s {
	var st s
	dir := config.Get().StorageDirectory
	file, err := os.ReadFile(dir + "/prv.pem")
	if err != nil {
		return &st
	}

	st.prv = privatekey.FromPem(string(file))
	st.pub = st.prv.PublicKey()
	st.addr = string(hasher.EncryptString(st.pub.ToString(false)))
	return &st
}

func Init() {
	var st s
	rand.Seed(time.Now().UnixNano())
	st.prv = privatekey.New(curve.Secp256k1, big.NewInt(rand.Int63()))
	st.save()
}

func Validate(data *block_data.ChainStored) bool {
	pubKey := publickey.FromString(data.PubKey, curve.Secp256k1, true)
	hash := hasher.EncryptString(fmt.Sprintf("%v", *data))
	return ecdsa.Verify(
		string(hash),
		signature.FromBase64(data.Sign),
		&pubKey,
	)
}

func Sign(data *block_data.ChainStored) {
	st := storage()
	hash := hasher.EncryptString(fmt.Sprintf("%v", *data))
	data.Sign = ecdsa.Sign(string(hash), &st.prv).ToBase64()
	data.PubKey = st.pub.ToString(false)
}

func GetAddr() string {
	return storage().addr
}
