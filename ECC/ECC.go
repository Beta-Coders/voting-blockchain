package ECC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"hash"
	"io"
	"os"
)

// GenKeys generates 2 pairs of keys using ECC
func GenKeys(data string) (*ecdsa.PrivateKey, ecdsa.PublicKey) {
	publicKeyCurve := elliptic.P256()
	privateKey := new(ecdsa.PrivateKey)
	privateKey, err := ecdsa.GenerateKey(publicKeyCurve, rand.Reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var publicKey ecdsa.PublicKey
	publicKey = privateKey.PublicKey
	//var h hash.Hash
	//h = md5.New()
	//_, err = io.WriteString(h, data)
	//if err != nil {
	//	return nil, ecdsa.PublicKey{}
	//}
	//signhash := h.Sum(nil)
	//signature, serr := ecdsa.SignASN1(rand.Reader, privateKey, signhash)
	//if serr != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//privatekey2, err := ecdsa.GenerateKey(publicKeyCurve, rand.Reader)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//var pubkey2 ecdsa.PublicKey
	//pubkey2 = privatekey2.PublicKey
	return privateKey, publicKey
}
func GenSign(data string, privateKey *ecdsa.PrivateKey) ([]byte, []byte) {
	var h hash.Hash
	h = md5.New()
	_, err := io.WriteString(h, data)
	if err != nil {
		return nil, nil
	}
	signhash := h.Sum(nil)
	signature, serr := ecdsa.SignASN1(rand.Reader, privateKey, signhash)
	if serr != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return signature, signhash
}

func GenPubKey(privateKey *ecdsa.PrivateKey) ecdsa.PublicKey {
	return privateKey.PublicKey
}
