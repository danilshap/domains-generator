package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/danilshap/domains-generator/pkg/config"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testStore Store
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := config.Load("../../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	testStore = NewStore(testDB)
	os.Exit(m.Run())
}
