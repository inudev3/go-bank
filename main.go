package main

import (
	"database/sql"
	"github.com/inudev5/go-bank/api"
	db "github.com/inudev5/go-bank/db/sqlc"
	"github.com/inudev5/go-bank/util"
	_ "github.com/lib/pq"
	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:32518458@localhost:5432/go_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBsource)
	if err != nil {
		log.Fatalln("cannot connect to db", err)
	}
	//testQueries = New(testDb)
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalln("cannot start server")
	}
}
