package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gobank/db/sqlc"
	"log"
	"os"
	testing "testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:32518458@localhost:5432/go_bank?sslmode=disable"
)

var testQueries *db.Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDb, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalln("cannot connect to db", err)
	}
	testQueries = db.New(testDb)
	os.Exit(m.Run()) //m.Run returns the test result, which is passed to os.Exit to report the result.
}
