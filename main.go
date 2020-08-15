package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mrzack99s/mrzack-dns-server/configs"

	"github.com/miekg/dns"
	except "github.com/mrzack99s/mrzack-dns-server/exceptions"
	"github.com/mrzack99s/mrzack-dns-server/optionals"
)

var domainsToAddresses map[string]string

var blockWords map[string][]string

type handler struct{}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
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
}

func main() {

	configs.ParseSystemConfig()
	domainsToAddresses = configs.ParseRecords()
	blockWords = configs.ParseBlockWords()

	srv := &dns.Server{Addr: ":" + strconv.Itoa(53), Net: "udp"}
	srv.Handler = &handler{}

	year, month, day := time.Now().Date()
	logFileName := "dns-log-" + strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day) + ".log"

	f, err := os.OpenFile(configs.SystemConfig.SConfig.LogPath+"/"+logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	log.SetOutput(f)

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Failed to set udp listener %s\n", err.Error())
	}
}
