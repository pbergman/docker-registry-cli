package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/pbergman/docker-registry-cli/account"
	"github.com/pbergman/docker-registry-cli/http"
	"github.com/pbergman/docker-registry-cli/logger"
)

type Manager struct {
	tokens map[string]*Token
	user   *account.User
}

func NewManager(user *account.User) *Manager {
	return &Manager{
		user:   user,
		tokens: make(map[string]*Token, 0),
	}
}

func (m *Manager) GetToken(c *http.AuthChallenge) (*Token, error) {

	hash := c.GetHash()

	if token, ok := m.tokens[hash]; ok {
		logger.Logger.Debug(fmt.Sprintf("Found token \"%x\"", hash))
		if token.IsValid() {
			return token, nil
		} else {
			logger.Logger.Debug(fmt.Sprintf("Found token \"%x\" was expired.", hash))
			delete(m.tokens, hash)
		}
	}

	logger.Logger.Debug("Requesting new token")

	response, err := http.Client.Do(c.GetRequest(m.user).Raw())

	logger.Logger.Debug(fmt.Sprintf("New token \"%x\"", hash))
	logger.Logger.Debug("Status: " + response.Status)

	switch response.StatusCode {
	case 401:
		return nil, errors.New("Authentication failed,")
	case 200:
		break
	default:
		return nil, errors.New("Failed to require token.")
	}

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	token := &Token{}
	err = json.Unmarshal(body, token)

	if err != nil {
		return nil, err
	}

	m.tokens[hash] = token

	return token, nil
}
