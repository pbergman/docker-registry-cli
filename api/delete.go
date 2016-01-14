package api

import (
	"fmt"
	"sync"

	"github.com/pbergman/docker-registry-cli/config"
	"github.com/pbergman/docker-registry-cli/http"
	"github.com/pbergman/docker-registry-cli/logger"
)

func Delete(repository, tag string, force bool) {
	manifest := GetManifest(repository, tag, true)

	if force {
		var wg sync.WaitGroup
		deleteBlob := func(digest string) {
			url := fmt.Sprintf("%s/v2/%s/blobs/%s", config.Config.RegistryHost, repository, digest)
			logger.Logger.Debug("Requesting DELETE for url: " + url)
			req := http.NewRequest(url)
			req.SetMethod("DELETE")
			resp, err := req.Do()
			logger.Logger.CheckError(err)
			resp, err = resolve(resp, false)
			logger.Logger.CheckWarning(err)
			if err != nil {
				logger.Logger.Debug("Removed: " + digest)
			}
			wg.Done()
		}

		for i, layer := range manifest.FsLayers {
			wg.Add(1)
			// first one not so token can be cached
			if i == 0 {
				deleteBlob(layer["blobSum"])
			} else {
				go deleteBlob(layer["blobSum"])
			}
		}
		wg.Wait()
	}

	url := fmt.Sprintf("%s/v2/%s/manifests/%s", config.Config.RegistryHost, repository, manifest.Digest)
	logger.Logger.Debug("Requesting DELETE for url: " + url)
	req := http.NewRequest(url)
	req.SetMethod("DELETE")
	resp, err := req.Do()
	resp, err = resolve(resp, true)
	logger.Logger.CheckError(err)
}
