package container

import (
	"database/sql"
	"fmt"
	"log"
	"votingblockchain/chain"
)

type Container interface {
	GetDB() *sql.DB
	GetBC() *chain.BlockChain
}
type container struct {
	db *sql.DB
	bc *chain.BlockChain
}

func NewContainer(db *sql.DB, bc *chain.BlockChain) Container {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected!")
	return &container{db: db, bc: bc}
}
func (t *container) GetDB() *sql.DB {
	return t.db
}
func (t *container) GetBC() *chain.BlockChain {
	return t.bc
}
