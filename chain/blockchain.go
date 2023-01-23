package chain

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"votingblockchain/block"
)

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

const BlocksBucket = "blocks"
const DbFile = "blockchain.db"

func (bc *BlockChain) AddBlock(candidateName string, username string, publicKey *ecdsa.PublicKey, signature []byte, signhash []byte) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	//_, pubkey, _, _, signature, signhash := ECC.GenKeys(username)
	newBlock := block.NewBlock(lastHash, candidateName, username, signature, signhash, *publicKey)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Fatal(err)
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})
}

func NewBlockChain() *BlockChain {
	var tip []byte
	db, err := bolt.Open(DbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

		if b == nil {
			genesis := block.NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(BlocksBucket))
			fmt.Println("creating new blockchain...")
			if err != nil {
				log.Fatal(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := BlockChain{tip, db}

	return &bc
}
func (bc *BlockChain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{CurrentHash: bc.tip, Db: bc.db}
	return bci
}
