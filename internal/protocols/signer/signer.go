package signer

import (
	"bvpn-prototype/internal/logger"
	"bvpn-prototype/internal/storage/config"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/curve"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/ecdsa"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/privatekey"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/publickey"
	"golang.org/x/crypto/sha3"
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
	st.addr = addrFromPub(st.pub.ToString(false))
	return &st
}

func addrFromPub(i string) string {
	hash := sha3.New256()
	hash.Write([]byte(i))
	return string(hash.Sum(nil))
}

func Init() {
	var st s
	rand.Seed(time.Now().UnixNano())
	st.prv = privatekey.New(curve.Secp256k1, big.NewInt(rand.Int63()))
	st.save()
}

func Sign(msg []byte) string {
	st := storage()
	signature := ecdsa.Sign(string(msg), &st.prv)
	return signature.ToBase64()
}

func GetAddr() string {
	return storage().addr
}

func GetPubKey() string {
	return storage().pub.ToString(false)
}
