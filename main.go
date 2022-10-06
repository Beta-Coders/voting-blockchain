package main

import (
	"fmt"
	"votingblockchain/chain"
	"votingblockchain/cli"
)

func main() {

	fmt.Println("starting up ....")
	bc := chain.NewBlockChain()
	defer bc.Iterator().Db.Close()

	cli := *&cli.CLI{Bc: bc}
	cli.Run()

	// for _, block := range bc.Block {
	// 	fmt.Printf("Prev. hash: %x\n", block.PrevHash)
	// 	fmt.Printf("Data: %s\n", block.Data)
	// 	fmt.Printf("Hash: %x\n", block.Hash)
	// 	x, y := elliptic.Unmarshal(elliptic.P256(), block.PubKey)
	// 	fmt.Printf("SignHash: %v\n", ecdsa.VerifyASN1(&ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, block.SignHash, block.Signature))
	// 	fmt.Println()
	// }
	// it := bc.Iterator()
	// for {

	// 	block := it.Next()

	// 	fmt.Printf("Prev. hash: %x\n", block.PrevHash)
	// 	fmt.Printf("Data: %s\n", block.Data)
	// 	fmt.Printf("Hash: %x\n", block.Hash)
	// 	x, y := elliptic.Unmarshal(elliptic.P256(), block.PubKey)
	// 	fmt.Printf("SignHash: %v\n", ecdsa.VerifyASN1(&ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, block.SignHash, block.Signature))
	// 	fmt.Println()

	// 	if len(block.PrevHash) == 0 {
	// 		break
	// 	}
	// }
}

// pair  1 admin - 2 pair private sign blocks pub

// check vote  || login || fake id || key unique || network verify || block data -> pub key || data hash
