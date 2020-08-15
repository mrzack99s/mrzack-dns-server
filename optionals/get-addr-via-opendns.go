package optionals

import (
	"context"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mrzack99s/mrzack-dns-server/configs"
)

var (
	ipRegex, _ = regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
)

func IsIpv4Regex(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")
	return ipRegex.MatchString(ipAddress)
}

func GetIPAddrFromOpenDNS(domain string) string {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, "udp", configs.SystemConfig.SConfig.ForwarderAddress)
		},
	}

	ips, err := r.LookupIPAddr(context.Background(), domain)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
	} else {
		for _, ip := range ips {
			if IsIpv4Regex(ip.String()) {
				return ip.String()
			}
		}
	}

	return "NOT_FOUND"
}
