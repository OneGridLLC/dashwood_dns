package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	lastAccess = "0"
	dnsdb      *sql.DB
	rwMutex    sync.RWMutex
)

var records = map[string][]string{} /*{
	"google.com.":    "64.233.177.102",
	"microsoft.com.": "20.84.181.62",
	"fake.service.":  "127.0.0.1",
}*/

func initDB() {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", *dbhost, *dbport, *dbuser, password, *dbname)

	// rand.Seed(int64(time.Now().UnixNano()))

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	//defer db.Close() does not work with coroutine

	dnsdb = db

	err = dnsdb.Ping()
	if err != nil {
		panic(err)
	}

	err = fetchRecords(true)
	if err != nil {
		panic(err)
	}

	go fetchRecordsRoutine() // records updating need not be synchronous
}

// Moved to controller
// func newLastAccess() string {
// 	return strconv.FormatInt(time.Now().UnixNano(), 36) + strconv.FormatUint(rand.Uint64(), 36)
// }

/* -- DATABASE STRUCTURE

Check controller/database_template.txt
got too big to copy into here

*/

const fetchRecordsSQL = `SELECT domain, address FROM records;`

// const testAccessMatchSQL = `SELECT value FROM pairs WHERE value = ?;`
const getLastAccessSQL = `SELECT value FROM pairs WHERE key = 'lastAccess'`

func testAccessMatch() (bool, string, error) {
	row := dnsdb.QueryRow(getLastAccessSQL)
	var identifier string
	err := row.Scan(&identifier)
	if err != nil {
		return false, "", err // any error
	}

	if identifier == lastAccess {
		return true, identifier, nil // value did not match
	}
	return false, identifier, nil // value did not match
}

type recordSmall struct { // not a full record but all that is necessary for this program
	domain  string
	address string
}

func fetchRecords(force bool) error {
	if !force {
		matches, newAccess, err := testAccessMatch()

		if err != nil {
			return err
		}

		if matches {
			return nil
		}

		lastAccess = newAccess // doesn't match, update it
	}

	rows, err := dnsdb.Query(fetchRecordsSQL)
	if err != nil {
		return err
	}

	data := []recordSmall{}

	for rows.Next() {
		i := recordSmall{}
		err := rows.Scan(&i.domain, &i.address)
		if err != nil {
			return err
		}
		data = append(data, i)
	}

	newRecords := map[string][]string{}

	for _, r := range data {
		newRecords[r.domain] = append(newRecords[r.domain], r.address)
		//newRecords[r.domain] = r.address
	}

	rwMutex.RLock() // prevent wacky antics
	records = newRecords
	rwMutex.RUnlock()

	fmt.Println("DB - Fetched new records")

	return nil

}

func fetchRecordsRoutine() {
	for range time.Tick(*dnsRefreshPeriod) {
		err := fetchRecords(false)
		if err != nil {
			fmt.Println("DB - Error fetching records: ", err)
		}
	}
}
