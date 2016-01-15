package config

import (
	"encoding/json"
	"github.com/pbergman/docker-registry-cli/account"
	"github.com/pbergman/docker-registry-cli/logger"
	"github.com/pbergman/docker-registry-cli/token"
	"os"
	"os/user"
)

type config struct {
	User         *account.User `json:"user"`
	TokenManager *token.Manager
	Verbose      bool   `json:"verbose,omitempty"`
	RegistryHost string `json:"registry-host"`
	Input        map[string]interface{}
	Command      int
}

var Config *config

func init() {

	logger.Logger.Pause(10)
	cnf, err := newConfig()
	Config = cnf
	Config.ParseInput()
	logger.SetHandler(Config.Verbose)
	logger.Logger.Resume()

	if err != nil {
		logger.Logger.Error(err)
		os.Exit(1)
	}
}

func newConfig() (*config, error) {

	newUser := account.NewEmptyUser()

	config := &config{
		User:         newUser,
		Input:        make(map[string]interface{}),
		TokenManager: token.NewManager(newUser),
	}

	user, err := user.Current()

	if err != nil {
		return config, err
	}

	file := user.HomeDir + "/.docker-registry/conf.json"

	if _, err := os.Stat(file); !os.IsNotExist(err) {

		logger.Logger.Debug("Found config file: " + file)
		fd, err := os.Open(file)

		if err == nil {
			parser := json.NewDecoder(fd)
			if err := parser.Decode(config); err != nil {
				return config, err
			} else {
				return config, nil
			}
		} else {
			return config, err
		}
	}

	return config, nil
}
