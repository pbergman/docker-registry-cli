package api

import (
	"github.com/pbergman/docker-registery-cli/config"
	"github.com/pbergman/docker-registery-cli/http"
	"github.com/pbergman/docker-registery-cli/logger"
)

// ApiCheck will check if the server implements the registry api v2
// If authentication (when needed) fails or server returns no 200
// or 401 status code it wil exit,
// see: https://docs.docker.com/registry/spec/api/#api-version-check
func ApiCheck() {
	logger.Logger.Debug("Server version check....")
	url := config.Config.RegistryHost + "/v2/"
	logger.Logger.Debug("Requesting HEAD for url: " + url)
	req := http.NewRequest(url)
	req.SetMethod("HEAD")
	req.AddHeader("Host", req.Raw().Host)
	resp, err := req.Do()
	logger.Logger.CheckError(err)
	logger.Logger.Debug("The registry implements the V2(.1) registry API.")
	if resp.StatusCode == 401 {
		logger.Logger.Debug("Server required authentication.")
	}
	resp, err = resolve(resp, true)
	if resp.StatusCode == 200 {
		logger.Logger.Debug("Authentication success.")
	}
}
