package expstore_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	expstore "github.com/gmaschi/log-exp-eval/internal/services/datastore/postgresql/exp"
	_ "github.com/lib/pq"
)

var testStore expstore.Store
var testDB *sql.DB

const (
	testHostDB     = "localhost"
	testDBPort     = "5434"
	testDBUser     = "myuser"
	testDBPassword = "mypassword"
	testDBName     = "eval-test"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			testHostDB, testDBPort, testDBUser, testDBPassword, testDBName),
	)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		os.Exit(1)
	}

	testStore = expstore.NewStore(testDB)

	os.Exit(m.Run())
}
