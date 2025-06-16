package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const dbSource = "postgresql://root:root@localhost:5432/ocr-database?sslmode=disable"

var testQueries *Queries
var testStore *Store
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	testDB, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	testStore = &Store{
		db:      testDB,
		Queries: testQueries,
	}

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

