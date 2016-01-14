package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pbergman/docker-registry-cli/config"
	"github.com/pbergman/docker-registry-cli/http"
	"github.com/pbergman/docker-registry-cli/logger"
)

type Repositories struct {
	Images []string `json:"repositories"`
}

// see: https://docs.docker.com/registry/spec/api/#listing-repositories
func GetRepositories() *Repositories {
	url := config.Config.RegistryHost + "/v2/_catalog"
	logger.Logger.Debug("Requesting url: " + url)
	req := http.NewRequest(url)
	resp, err := req.Do()
	logger.Logger.CheckError(err)
	resp, err = resolve(resp, true)
	body, err := ioutil.ReadAll(resp.Body)
	logger.Logger.CheckError(err)
	repository := &Repositories{}
	err = json.Unmarshal(body, repository)
	logger.Logger.CheckError(err)
	return repository
}
