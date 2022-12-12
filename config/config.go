package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config interface {
	GetHost() string
	GetDBName() string
	GetUser() string
	GetPassword() string
	GetPort() int
}

type config struct {
	host     string
	dbName   string
	user     string
	password string
	port     int
}

// Load configuration with -env default : develop
func Load() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	host := os.Getenv("HOST")
	user := os.Getenv("USER")
	dbName := os.Getenv("DB_NAME")
	password := os.Getenv("PASSWORD")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	return &config{host: host, dbName: dbName, user: user, password: password, port: port}
}

func (t *config) GetHost() string {
	return t.host
}
func (t *config) GetDBName() string {
	return t.dbName
}
func (t *config) GetUser() string {
	return t.user
}
func (t *config) GetPassword() string {
	return t.password
}
func (t *config) GetPort() int {
	return t.port
}
