package main

import "time"

const (
	dnsPort          = 53532
	dbport           = 5432
	dbhost           = "localhost"
	dbuser           = "postgres"
	dbname           = "dns"
	dnsRefreshPeriod = 5 * time.Minute
)

func main() {

	initDB()
	initDNS()

}
