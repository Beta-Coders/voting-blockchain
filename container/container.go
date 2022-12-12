package container

import (
	"database/sql"
	"fmt"
	"log"
)

type Container interface {
	GetDB() *sql.DB
}
type container struct {
	db *sql.DB
}

func NewContainer(db *sql.DB) Container {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected!")
	return &container{db: db}
}
func (t *container) GetDB() *sql.DB {
	return t.db
}
