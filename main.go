package main

import (
	"crypto/ecdsa"
	"fmt"
	"votingblockchain/chain"
)

func main() {

	fmt.Println("starting up ....")
	bc := chain.NewBlockChain()
	bc.AddBlock("TEST DATA")
	bc.AddBlock("TEST DATA 2")
	for _, block := range bc.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("SignHash: %v\n", ecdsa.VerifyASN1(&block.PubKey, block.SignHash, block.Signature))
		fmt.Println()
	}
}

// pair  1 admin - 2 pair private sign blocks pub

// check vote  || login || fake id || key unique || network verify || block data -> pub key || data hash