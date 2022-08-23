package main

const dnsPort = 53532
const apiPort = ":43023"

func main() {

	initDB()
	initDNS()
	webHandlers()

}
