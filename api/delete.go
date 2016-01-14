package api

import (
	"fmt"
	"sync"

	"github.com/pbergman/docker-registry-cli/config"
	"github.com/pbergman/docker-registry-cli/http"
	"github.com/pbergman/docker-registry-cli/logger"
)

type blobList []string

func (b *blobList) has(blob string) bool {
	for _, n := range *b {
		if n == blob {
			return true
		}
	}
	return false
}

func (b *blobList) add(blob string) {
	if false == b.has(blob) {
		*b = append(*b, blob)
	}
}

func Delete(repository, tag string, force bool) {

	var wg sync.WaitGroup
	blobs := new(blobList)
	for _, repos := range *GetList() {
		for _, tagName := range repos.Tags {
			if repos.Name == repository && tagName == tag {
				continue
			}
			wg.Add(1)
			go func(list *blobList, repos, tag string) {
				manifests := GetManifest(repos, tag, true)
				for _, blob := range manifests.FsLayers {
					list.add(blob["blobSum"])
				}
				wg.Done()
			}(blobs, repos.Name, tagName)
		}
	}
	wg.Wait()
	manifest := GetManifest(repository, tag, true)
	remove := new(blobList)
	for _, layer := range manifest.FsLayers {
		if false == blobs.has(layer["blobSum"]) {
			logger.Logger.Info(fmt.Sprintf("Adding layer %s to queue for removal", layer["blobSum"]))
			remove.add(layer["blobSum"])
		} else {
			logger.Logger.Info(fmt.Sprintf("Skipping layer %s (is used by other images)", layer["blobSum"]))
		}
	}

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

	for i, layer := range *remove {
		wg.Add(1)
		// first one not so token can be cached
		if i == 0 {
			deleteBlob(layer)
		} else {
			go deleteBlob(layer)
		}
	}

	url := fmt.Sprintf("%s/v2/%s/manifests/%s", config.Config.RegistryHost, repository, manifest.Digest)
	logger.Logger.Debug("Requesting DELETE for url: " + url)
	req := http.NewRequest(url)
	req.SetMethod("DELETE")
	resp, err := req.Do()
	resp, err = resolve(resp, true)
	logger.Logger.CheckError(err)
}
