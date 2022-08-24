package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

const (
	dnsPort          = 53
	dbport           = 5432
	dbhost           = "localhost"
	dbuser           = "postgres"
	dbname           = "dns"
	dnsRefreshPeriod = 5 * time.Second
)

func main() {

	initDB()
	initDNS()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt) // I don't think this actually works, and I don't really think it's necessary, but I'll keep it for now.

	go func() {
		sig := <-c
		fmt.Printf("Signal: %s | Exiting...\n", sig)
		dnsdb.Close()
		os.Exit(1)
	}()

}
