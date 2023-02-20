package ovpn

import (
	"bufio"
	"fmt"
	"github.com/forestsay/oar-dns/util"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type ClientRecord struct {
	Address     string
	RefreshTime time.Time
}

type ManagementClient struct {
    EndPoint string

    routeMap sync.Map
}

func (m *ManagementClient) clientReader(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, _ := reader.ReadString('\n')
		m.onReceiveMessage(message)
	}
}

func (m *ManagementClient) onReceiveMessage(message string) {
	if strings.HasPrefix(message, "ROUTING_TABLE") {
		arr := strings.Split(message, "\t")
		log.Println("Added:", arr[2], arr[1])
		m.routeMap.Store(arr[2], ClientRecord{arr[1], time.Now()})
	}
}

func (m *ManagementClient) clientWriter(conn net.Conn) {
	_, err := fmt.Fprintf(conn, util.GetEnvVar("OPENVPN-MANAGEMENT-INTERFACE-PASSWORD", "password")+"\r\n")
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

func (m *ManagementClient) GuessDomain(domain string) (result ClientRecord, ok bool) {
	mapResult, ok := m.routeMap.Load(domain)
	if ok {
		result = mapResult.(ClientRecord)
		return
	}

	domain = strings.TrimSuffix(domain, ".")
	mapResult, ok = m.routeMap.Load(domain)
	if ok {
		result = mapResult.(ClientRecord)
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
	mapResult, ok = m.routeMap.Load(s)
	if ok {
		result = mapResult.(ClientRecord)
		return
	}

	ok = false
	return
}

func (m *ManagementClient) StartOpenVPNClient() {
	c, err := net.Dial("tcp", m.EndPoint)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	go m.clientReader(c)
	go m.clientWriter(c)
}
