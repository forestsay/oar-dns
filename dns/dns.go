package dns

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/forestsay/oar-dns/ovpn"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
	"time"
)

type Server struct {
	OpenVPNManagementClient *ovpn.ManagementClient
	Apollo                  agollo.Client
}

func (d *Server) resolveARecordFromOVPN(domain string) (result ovpn.ClientRecord, ok bool) {
	if !strings.HasSuffix(domain, "oar.moe.") {
		ok = false
		return
	}
	domain = domain[:len(domain)-len("oar.moe.")]

	result, ok = d.OpenVPNManagementClient.GuessDomain(domain)
	if !ok {
		ok = false
		return
	}
	validAfter := time.Now().Add(-time.Second * 60)
	if !result.RefreshTime.After(validAfter) {
		ok = false
		return
	}

	ok = true
	return
}

func (d *Server) resolveFromApollo(t string, domain string) (string, bool) {
	result := d.Apollo.GetStringValue(fmt.Sprintf("record:%s:%s", t, domain), "")
	if result != "" {
		return result, true
	}
	return "", false
}

func (d *Server) resolveDomain(domain string) (result string, ok bool) {
	apolloResult, ok := d.resolveFromApollo("A", domain)
	if ok {
		result = apolloResult
	}

	ovpnResult, ok := d.resolveARecordFromOVPN(domain)
	if ok {
		result = ovpnResult.Address
	}

	return "", false
}

func (d *Server) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		result, ok := d.resolveDomain(domain)
		if ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(result),
			})
		} else {
			msg.Rcode = dns.RcodeNameError
		}
	default:
		msg.Rcode = dns.RcodeNameError
	}
	_ = w.WriteMsg(&msg)
}

type ServerConfig struct {
	OpenVPNManagementClient *ovpn.ManagementClient
	Listen                  string
	Apollo                  agollo.Client
}

func (d *ServerConfig) StartDNSServer() {
	srv := &dns.Server{
		Addr: d.Listen,
		Net:  "udp",
	}
	srv.Handler = &Server{
		OpenVPNManagementClient: d.OpenVPNManagementClient,
		Apollo:                  d.Apollo,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
}
