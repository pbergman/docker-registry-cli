package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pbergman/docker-registery-cli/config"
	"github.com/pbergman/docker-registery-cli/http"
	"github.com/pbergman/docker-registery-cli/logger"
)

type Tag struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func GetTags(repository string) *Tag {
	url := fmt.Sprintf("%s/v2/%s/tags/list", config.Config.RegistryHost, repository)
	logger.Logger.Debug("Requesting url: " + url)
	req := http.NewRequest(url)
	resp, err := req.Do()
	logger.Logger.CheckError(err)
	resp, err = resolve(resp, true)
	body, err := ioutil.ReadAll(resp.Body)
	logger.Logger.CheckError(err)
	tag := &Tag{}
	err = json.Unmarshal(body, tag)
	logger.Logger.CheckError(err)

	return tag
}
