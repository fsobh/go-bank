package db

import (
	"database/sql"
	"github.com/fsobh/simplebank/util"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

// Global variables in testing
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	//load env file
	config, err := util.LoadConfig("../..") // pass in the path relative to this file

	if err != nil {
		log.Fatal("Cannot load env variables", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Can not connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
