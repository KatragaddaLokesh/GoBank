package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var conn *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:Kashira@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	conn, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot Connect To database:", err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())

}
