package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"net/http"
	_ "net/http/pprof"
)

const (
	profileServe = ":51280" // port that probably isn't used
)

var (
	dbhost   = flag.String("dbhost", "localhost", "Postgres address")
	dbport   = flag.Int("dbport", 5432, "Postgres port")
	dbuser   = flag.String("dbuser", "postgres", "Postgres username")
	dbname   = flag.String("dbname", "dns", "Postgres database name")
	profile  = flag.Bool("profile", false, "Enable profiling")
	password = os.Getenv("PostgresPass")

	dnsPort          = flag.Int("dnsport", 53, "Port to serve DNS")
	dnsRefreshPeriod = flag.Duration("dnsrefreshperiod", 1*time.Minute, "") // arbitrary amount of time to check for DNS changes
)

func main() {

	flag.Parse()

	initDB()
	defer initDNS() // blocks main thread

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt) // I don't really think this is necessary, but I'll keep it for now.

	go func() {
		sig := <-c
		fmt.Printf("Signal: %s | Exiting...\n", sig)
		dnsdb.Close()
		os.Exit(1)
	}()

	// profiling
	if *profile {
		go func() {
			http.ListenAndServe(profileServe, nil)
		}()
	}
}
