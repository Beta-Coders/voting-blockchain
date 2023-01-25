package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Hash          []byte
	PrevHash      []byte
	CandidateName []byte
	Signature     []byte
	SignHash      []byte
	PubKey        []byte
	Username      []byte
}

// NewBlock creates a new block in the blockchain with the given data and previous hash
func NewBlock(prevHash []byte, candidate string, username string, signature []byte, signhash []byte, pubKey ecdsa.PublicKey) *Block {
	var pubkey []byte
	var pubkeyEmpty ecdsa.PublicKey
	if pubKey != pubkeyEmpty {
		pubkey = elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	}
	b := &Block{time.Now().Unix(), []byte{}, prevHash, []byte(candidate), signature, signhash, pubkey, []byte(username)}
	b.setHash()
	return b
}

func (b *Block) setHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	msg := bytes.Join([][]byte{b.PrevHash, b.Username, timestamp}, []byte{})
	hash := sha256.Sum256(msg)
	b.Hash = hash[:]
}

// NewGenesisBlock creates the genesis block when the blockchain is created
func NewGenesisBlock() *Block {
	var pubkey ecdsa.PublicKey
	return NewBlock([]byte{}, "Genesis Block", "Genesis Block", []byte{}, []byte{}, pubkey)
}
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Fatal(err)
	}

	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Fatal(err)
	}

	return &block
}

// pub priv gen ->
