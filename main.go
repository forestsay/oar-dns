package main

import (
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/forestsay/oar-dns/apollo"
	"github.com/forestsay/oar-dns/dns"
	"github.com/forestsay/oar-dns/log"
	"github.com/forestsay/oar-dns/ovpn"
	"github.com/forestsay/oar-dns/util"
)

var GitCommitId string

func main() {
	log.Logger.Print("GitCommitId: " + GitCommitId)

	omc := &ovpn.ManagementClient{
		EndPoint: util.GetEnvVar("OPENVPN-MANAGEMENT-INTERFACE-ENDPOINT", "127.0.0.1:27273"),
	}
	omc.StartOpenVPNClient()

	apolloConfig := &config.AppConfig{
		AppID:             util.GetEnvVar("OARDNS-APOLLO-APPID", "oardns"),
		Cluster:           util.GetEnvVar("OARDNS-APOLLO-CLUSTER", "dev"),
		IP:                util.GetEnvVar("OARDNS-APOLLO-IP", "http://127.0.0.1:8080"),
		NamespaceName:     util.GetEnvVar("OARDNS-APOLLO-NAMESPACE", "application"),
		IsBackupConfig:    util.GetEnvVar("OARDNS-APOLLO-ISBACKUPCONFIG", "true") == "true",
		Secret:            util.GetEnvVar("OARDNS-APOLLO-SECRET", ""),
		SyncServerTimeout: 10,
	}
	cfg := apollo.StartApollo(apolloConfig)

	ns := &dns.ServerConfig{
		OpenVPNManagementClient: omc,
		Apollo:                  cfg,
		Listen:                  util.GetEnvVar("DNS-SERVER-LISTEN", ":53"),
	}
	ns.StartDNSServer()

	select {}
}
