package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	password   = os.Getenv("PostgresPass")
	lastAccess = ""
	dnsdb      *sql.DB
	rwMutex    sync.RWMutex
)

var records = map[string]string{} /*{
	"google.com.":    "64.233.177.102",
	"microsoft.com.": "20.84.181.62",
	"fake.service.":  "127.0.0.1",
}*/

func initDB() {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, password, dbname)

	// rand.Seed(int64(time.Now().UnixNano()))

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	dnsdb = db

	err = dnsdb.Ping()
	if err != nil {
		panic(err)
	}

	err = fetchRecords(true)
	if err != nil {
		panic(err)
	}

	go fetchRecordsRoutine() // async records updating
}

// Moved to controller
// func newLastAccess() string {
// 	return strconv.FormatInt(time.Now().UnixNano(), 36) + strconv.FormatUint(rand.Uint64(), 36)
// }

/* -- DATABASE STRUCTURE

CREATE TABLE IF NOT EXISTS pairs (
	key TEXT NOT NULL,
	value TEXT NOT NULL,
	PRIMARY KEY (key)
); -- this contains the update match string

CREATE TABLE IF NOT EXISTS records (
	id uuid DEFAULT uuid_generate_v4 (),
	domain TEXT NOT NULL,
	address TEXT NOT NULL,
	created TIMESTAMP NOT NULL,
	PRIMARY KEY (domain)
); -- has to have an ID because a domain can link to more than one address

*/

const fetchRecordsSQL = `SELECT domain, address FROM records;`
const testAccessMatchSQL = `SELECT value FROM records WHERE value=?`

func testAccessMatch() (bool, string, error) {
	row := dnsdb.QueryRow(testAccessMatchSQL, lastAccess)
	var identifier string
	err := row.Scan(&identifier)
	if err != nil && err == sql.ErrNoRows {
		return false, "", nil // value did not match
	} else if err != nil {
		return false, "", err // arbitrary error
	}

	return true, identifier, nil // value did match
}

type recordSmall struct { // not a full record but all that is necessary for this program
	domain  string
	address string
}

func fetchRecords(force bool) error {
	if !force {
		matches, newAccess, err := testAccessMatch()

		lastAccess = newAccess
		if err != nil {
			return err
		}
		if matches {
			return nil
		}
	}

	rows, err := dnsdb.Query(fetchRecordsSQL)
	if err != nil {
		return err
	}

	data := []recordSmall{}

	for rows.Next() {
		i := recordSmall{}
		err := rows.Scan(i.domain, i.address)
		if err != nil {
			return err
		}
		data = append(data, i)
	}

	newRecords := map[string]string{}

	for _, r := range data {
		newRecords[r.domain] = r.address
	}

	rwMutex.RLock() // prevent wacky antics
	records = newRecords
	rwMutex.RUnlock()

	fmt.Println("DB - Fetched new records")

	return nil

}

func fetchRecordsRoutine() {
	for {
		time.Sleep(dnsRefreshPeriod)

		err := fetchRecords(false)
		if err != nil {
			fmt.Println("DB - Error fetching records: ", err)
		}
	}
}