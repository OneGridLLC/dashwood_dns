package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const (
	port   = 5432
	user   = "postgres"
	dbname = "dns"
)

var (
	host       = os.Getenv("PostgresHost")
	password   = os.Getenv("PostgresPass")
	lastAccess = newLastAccess()
)

var records = map[string]string{
	"google.com.":    "64.233.177.102",
	"microsoft.com.": "20.84.181.62",
	"fake.service.":  "127.0.0.1",
}

func initDB() {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	rand.Seed(int64(time.Now().UnixNano()))

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	fmt.Println(lastAccess)
}

func newLastAccess() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36) + strconv.FormatUint(rand.Uint64(), 36)
}

func fetchRecords() {

}
