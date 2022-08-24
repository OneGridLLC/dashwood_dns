package main

//ripped from https://gist.github.com/walm/0d67b4fb2d5daf3edd4fad3e13b162cb

import (
	"fmt"
	"log"
	"strconv"

	"github.com/miekg/dns"
)

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("DNS - Query for %s\n", q.Name)
			rwMutex.Lock() // prevent wacky antics
			ip := records[q.Name]
			rwMutex.Unlock()
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func initDNS() {
	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	// start server
	port := dnsPort
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("DNS - Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("DNS - Failed to start server: %s\n ", err.Error())
	}
}
