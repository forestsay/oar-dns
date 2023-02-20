package apollo

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/protocol/auth"
	"github.com/forestsay/oar-dns/log"
)

func StartApollo(c *config.AppConfig) agollo.Client {
	agollo.SetLogger(log.Logger)

	//extension.SetHTTPAuth(&AuthSignature{
	//	underlay: &sign.AuthSignature{},
	//})
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		log.Logger.Panic(err)
	}

	return client
}

type AuthSignature struct {
	underlay auth.HTTPAuth
}

func (a AuthSignature) HTTPHeaders(url string, appID string, secret string) map[string][]string {
	if secret == "" {
		return nil
	}
	return a.underlay.HTTPHeaders(url, appID, secret)
}
