package db

import (
	"database/sql"
	"github.com/inudev5/go-bank/util"
	_ "github.com/lib/pq"
	"log"
	"os"
	testing "testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:32518458@localhost:5432/go_bank?sslmode=disable"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal(err)
	}
	testDb, err = sql.Open(config.DBDriver, config.DBsource)
	if err != nil {
		log.Fatalln("cannot connect to db", err)
	}
	testQueries = New(testDb)
	os.Exit(m.Run()) //m.Run returns the test result, which is passed to os.Exit to report the result.
}
