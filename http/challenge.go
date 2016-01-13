package http

import (
	"crypto/sha1"
	"github.com/pbergman/docker-registery-cli/account"
	"github.com/pbergman/docker-registery-cli/logger"
	"io"
)

type AuthChallenge struct {
	Service string `header:"service"`
	Realm   string `header:"realm"`
	Scope   string `header:"scope"`
	Error   string `header:"error"`
	Scheme  string
}

func (a *AuthChallenge) GetHash() string {
	hash := sha1.New()
	io.WriteString(hash, a.Service)
	io.WriteString(hash, a.Realm)
	io.WriteString(hash, a.Scope)
	return string(hash.Sum(nil))
}

func (a *AuthChallenge) GetRequest(u *account.User) *Request {

	request := NewRequest(a.Realm)
	request.SetBasicAuth(u.Username, u.Password)

	logger.Logger.Debug("realm: " + a.Realm)
	logger.Logger.Debug("service: " + a.Service)

	query := request.Raw().URL.Query()
	query.Set("service", a.Service)

	if a.Scope != "" {
		query.Set("scope", a.Scope)
		logger.Logger.Debug("scope: " + a.Scope)
	}

	request.Raw().URL.RawQuery = query.Encode()

	return request
}
