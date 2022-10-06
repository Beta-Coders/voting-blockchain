package chain

import (
	"log"
	"votingblockchain/block"

	"github.com/boltdb/bolt"
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
	i.CurrentHash = deserializedBlock.PrevHash
	return deserializedBlock
}
