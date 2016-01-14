package api

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/pbergman/docker-registry-cli/config"
	"github.com/pbergman/docker-registry-cli/http"
	"github.com/pbergman/docker-registry-cli/logger"
)

func GetSize(repository, tag string) int {
	manifest := GetManifest(repository, tag, true)
	var wg sync.WaitGroup
	var size int
	for _, layer := range manifest.FsLayers {
		wg.Add(1)
		go func(digest string) {
			url := fmt.Sprintf("%s/v2/%s/blobs/%s", config.Config.RegistryHost, repository, digest)
			logger.Logger.Debug("Requesting HEAD for url: " + url)
			req := http.NewRequest(url)
			req.SetMethod("HEAD")
			resp, err := req.Do()
			logger.Logger.CheckError(err)
			resp, err = resolve(resp, true)
			logger.Logger.CheckError(err)
			// resp.Header.Get("Accept-Ranges")
			length, err := strconv.Atoi(resp.Header.Get("Content-Length"))
			logger.Logger.CheckError(err)
			size += length
			wg.Done()
		}(layer["blobSum"])
	}
	wg.Wait()
	return size
}
