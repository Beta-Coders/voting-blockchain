package cli

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"flag"
	"fmt"
	"log"
	"os"
	"votingblockchain/chain"
)

type CLI struct {
	Bc *chain.BlockChain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  blockjoden -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  chaindikhaen - print all the blocks of the blockchain")
}
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("blockjoden", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("chaindikhaen", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "blockjoden":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}

	case "chaindikhaen":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)

	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			cli.printUsage()
			os.Exit(1)
		}
		//cli.Bc.AddBlock(*addBlockData)
	}
	if printChainCmd.Parsed() {
		it := cli.Bc.Iterator()
		for {

			block := it.Next()

			fmt.Printf("Prev. hash: %x\n", block.PrevHash)
			//fmt.Printf("Data: %s\n", block.Data)
			fmt.Printf("Hash: %x\n", block.Hash)
			x, y := elliptic.Unmarshal(elliptic.P256(), block.PubKey)
			fmt.Printf("SignHash: %v\n", ecdsa.VerifyASN1(&ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, block.SignHash, block.Signature))
			fmt.Println()

			if len(block.PrevHash) == 0 {
				break
			}
		}
	}

}
