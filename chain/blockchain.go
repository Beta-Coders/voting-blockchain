package chain

import (
	"votingblockchain/ECC"
	"votingblockchain/block"
)

type BlockChain struct {
	Blocks []*block.Block
}

func (bc *BlockChain) AddBlock(data string) {
	_, pubkey, _, _, signature, signhash := ECC.GenKeys()
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := block.NewBlock(prevBlock.Hash, data, signature, signhash, pubkey)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*block.Block{block.NewGenesisBlock()}}
}
