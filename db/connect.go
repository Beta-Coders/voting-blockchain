package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func Connect(host string, user string, dbName string, password string, port int) *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	result, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

