package chain

import (
	"github.com/boltdb/bolt"
	"log"
	"votingblockchain/block"
)

type BlockchainIterator struct {
	CurrentHash []byte
	Db          *bolt.DB
}

func (i *BlockchainIterator) Next() *block.Block {
	var deserializedBlock *block.Block

	err := i.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		serializedBlock := b.Get(i.CurrentHash)
		deserializedBlock = block.DeserializeBlock(serializedBlock)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return deserializedBlock
}
