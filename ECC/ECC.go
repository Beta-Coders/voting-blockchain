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

func GenKeys() (*ecdsa.PrivateKey, ecdsa.PublicKey, *ecdsa.PrivateKey, ecdsa.PublicKey, []byte, []byte) {
	pubkeyCurve := elliptic.P256()
	privatekey := new(ecdsa.PrivateKey)
	privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var pubkey ecdsa.PublicKey
	pubkey = privatekey.PublicKey
	var h hash.Hash
	h = md5.New()
	io.WriteString(h, "This is a message to be signed and verified by ECDSA!")
	signhash := h.Sum(nil)
	signature, serr := ecdsa.SignASN1(rand.Reader, privatekey, signhash)
	if serr != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	privatekey2, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var pubkey2 ecdsa.PublicKey
	pubkey2 = privatekey2.PublicKey
	return privatekey, pubkey, privatekey2, pubkey2, signature, signhash
}
