package main

import (
	"bufio"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var GitCommitId string

type routeRecord struct {
	address     string
	refreshTime time.Time
}

var routeMap sync.Map

func clientReader(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, _ := reader.ReadString('\n')
		//fmt.Print("->: " + message)
		onReceiveMessage(message)
	}
}

func onReceiveMessage(message string) {
	if strings.HasPrefix(message, "ROUTING_TABLE") {
		arr := strings.Split(message, "\t")
		log.Println("Added:", arr[2], arr[1])
		routeMap.Store(arr[2], routeRecord{arr[1], time.Now()})
	}
}

func clientWriter(conn net.Conn) {
	_, err := fmt.Fprintf(conn, GetEnvVar("OPENVPN-MANAGEMENT-INTERFACE-PASSWORD", "password")+"\r\n")
	if err != nil {
		os.Exit(3)
	}
	_, err = fmt.Fprintf(conn, "status 3\r\n")
	if err != nil {
		os.Exit(3)
	}

	ticker := time.NewTicker(20 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				_, err = fmt.Fprintf(conn, "status 3\r\n")
				if err != nil {
					os.Exit(3)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func guessDomain(domain string) (result routeRecord, ok bool) {
	mapResult, ok := routeMap.Load(domain)
	if ok {
		result = mapResult.(routeRecord)
		return
	}

	domain = strings.TrimSuffix(domain, ".")
	mapResult, ok = routeMap.Load(domain)
	if ok {
		result = mapResult.(routeRecord)
		return
	}

	domainStep := strings.Split(domain, ".")
	s := ""
	for i, element := range domainStep {
		if i == 0 {
			s = element
		} else {
			s = element + "-" + s
		}
	}
	mapResult, ok = routeMap.Load(s)
	if ok {
		result = mapResult.(routeRecord)
		return
	}

	ok = false
	return
}

func parseDomain(domain string) (result routeRecord, ok bool) {
	result, ok = guessDomain(domain)
	if !ok {
		ok = false
		return
	}
	validAfter := time.Now().Add(-time.Second * 60)
	if !result.refreshTime.After(validAfter) {
		ok = false
		return
	}

	ok = true
	return
}

type dnsServer struct{}

func (dnsSrv *dnsServer) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		result, ok := parseDomain(domain)
		if ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(result.address),
			})
		} else {
			msg.Rcode = dns.RcodeNameError
		}
	default:
		msg.Rcode = dns.RcodeNameError
	}
	_ = w.WriteMsg(&msg)
}

func main() {
	log.Println("GitCommitId: " + GitCommitId)

	c, err := net.Dial("tcp", GetEnvVar("OPENVPN-MANAGEMENT-INTERFACE-ENDPOINT", "127.0.0.1:27273"))
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	go clientReader(c)
	go clientWriter(c)

	srv := &dns.Server{Addr: GetEnvVar("DNS-SERVER-LISTEN", ":53"), Net: "udp"}
	srv.Handler = &dnsServer{}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}

	select {}
}
