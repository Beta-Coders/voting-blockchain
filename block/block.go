package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp int64
	Hash      []byte
	PrevHash  []byte
	Data      []byte
	Signature []byte
	SignHash  []byte
	PubKey    ecdsa.PublicKey
}

func NewBlock(prevHash []byte, data string, signature []byte, signhash []byte, pubKey ecdsa.PublicKey) *Block {
	b := &Block{time.Now().Unix(), []byte{}, prevHash, []byte(data), signature, signhash, pubKey}
	b.setHash()
	return b
}

func (b *Block) setHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	msg := bytes.Join([][]byte{b.PrevHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(msg)
	b.Hash = hash[:]
}

func NewGenesisBlock() *Block {
	var pubkey ecdsa.PublicKey
	return NewBlock([]byte{}, "Genesis Block", []byte{}, []byte{}, pubkey)
}

// pub priv gen ->
