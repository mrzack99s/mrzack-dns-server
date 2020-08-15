package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/miekg/dns"
	"github.com/mrzack99s/mrzack-dns-server/configs"
	"github.com/mrzack99s/mrzack-dns-server/runtimes"
)

var domainsToAddresses map[string]string

var blockWords map[string][]string

type handler struct{}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	runtimes.Run(w, r, domainsToAddresses, blockWords)
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
