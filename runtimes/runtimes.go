package runtimes

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	except "github.com/mrzack99s/mrzack-dns-server/exceptions"

	"github.com/miekg/dns"
	"github.com/mrzack99s/mrzack-dns-server/optionals"
)

func runProcess(w dns.ResponseWriter, r *dns.Msg, domainsToAddresses map[string]string, blockWords map[string][]string) {

	except.Block{
		Try: func() {
			msg := dns.Msg{}
			msg.SetReply(r)
			switch r.Question[0].Qtype {
			case dns.TypeA:
				status := "SUCCESS"
				msg.Authoritative = true
				domain := msg.Question[0].Name

				except.Block{
					Try: func() {

						foundWordBlock := false
						for _, word := range blockWords["words"] {
							found := strings.Contains(domain, word)
							if found {
								foundWordBlock = true
							}
						}

						if !foundWordBlock {
							address, ok := domainsToAddresses[domain]
							if ok {
								msg.Answer = append(msg.Answer, &dns.A{
									Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
									A:   net.ParseIP(address),
								})
							} else {
								ipAddr := optionals.GetIPAddrFromOpenDNS(domain)
								if ipAddr != "NOT_FOUND" {
									msg.Answer = append(msg.Answer, &dns.A{
										Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
										A:   net.ParseIP(ipAddr),
									})
								} else {
									status = "FAILED"
								}
							}
						} else {
							address, ok := domainsToAddresses["block."]
							if ok {
								msg.Answer = append(msg.Answer, &dns.A{
									Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
									A:   net.ParseIP(address),
								})
								status = "HAS_BLOCKED"
							}
						}

						log.Printf("%s \t %s \t| STATUS | %s", w.RemoteAddr().String(), msg.Question[0].Name, status)
					},
					Catch: func(e except.Exception) {
						fmt.Printf("Caught %v\n", e)
					},
				}.Do()

			}
			w.WriteMsg(&msg)
		},
		Catch: func(e except.Exception) {
			fmt.Printf("Caught %v\n", e)
		},
	}.Do()

}

var (
	wg sync.WaitGroup
)

func Run(w dns.ResponseWriter, r *dns.Msg, domainsToAddresses map[string]string, blockWords map[string][]string) {
	go runProcess(w, r, domainsToAddresses, blockWords)
}
