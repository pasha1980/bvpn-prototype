package vpn_crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"golang.org/x/crypto/sha3"
)

/*
  Encoding Algo: RSA / SHA3-256
*/

func GeneratePair() ([]byte, []byte) {
	prv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := prv.PublicKey

	encodedPrv := x509.MarshalPKCS1PrivateKey(prv)
	encodedPub := x509.MarshalPKCS1PublicKey(&pub)

	return encodedPrv, encodedPub
}

func Decode(encrypted []byte, key []byte) ([]byte, error) {
	prv, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}

	decryptedBytes, err := prv.Decrypt(nil, encrypted, &rsa.OAEPOptions{Hash: crypto.SHA3_256})
	if err != nil {
		return nil, err
	}

	return decryptedBytes, nil
}

func Encode(plain []byte, key []byte) ([]byte, error) {
	pub, err := x509.ParsePKCS1PublicKey(key)
	if err != nil {
		return nil, err
	}

	encryptedBytes, err := rsa.EncryptOAEP(
		sha3.New256(),
		rand.Reader,
		pub,
		plain,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return encryptedBytes, nil
}
