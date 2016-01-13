package api

import (
	"fmt"
	"io/ioutil"

	"github.com/pbergman/docker-registery-cli/config"
	"github.com/pbergman/docker-registery-cli/http"
	"github.com/pbergman/docker-registery-cli/logger"
)

func Delete(repository, tag string) {
	manifest := GetManifest(repository, tag)
	url := fmt.Sprintf("%s/v2/%s/manifests/%s", config.Config.RegistryHost, repository, manifest.Digest)
	logger.Logger.Debug("Requesting DELETE for url: " + url)
	req := http.NewRequest(url)
	req.SetMethod("DELETE")
	resp, err := req.Do()
	resp, err = resolve(resp, true)
	logger.Logger.CheckError(err)
	_, err = ioutil.ReadAll(resp.Body)
	logger.Logger.CheckError(err)
}
