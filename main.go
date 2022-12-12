package main

import (
	"database/sql"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"votingblockchain/chain"
	"votingblockchain/config"
	"votingblockchain/container"
	"votingblockchain/db"
	"votingblockchain/router"
)

func main() {

	fmt.Println()
	fmt.Println("starting up ....")
	fmt.Println()
	bc := chain.NewBlockChain()
	defer func(Db *bolt.DB) {
		err := Db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(bc.Iterator().Db)
	conf := config.Load()
	newContainer := container.NewContainer(db.Connect(conf.GetHost(), conf.GetUser(), conf.GetDBName(), conf.GetPassword(), conf.GetPort()))
	defer func(Db *sql.DB) {
		err := Db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(newContainer.GetDB())
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
	router.Init(e, newContainer)
	e.Logger.Fatal(e.Start(":1323"))
	//app := *&cli.CLI{Bc: bc}
	//app.Run()
}

// pair  1 admin - 2 pair private sign blocks pub

// check vote  || login || fake id || key unique || network verify || block data -> pub key || data hash
