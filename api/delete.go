package api

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pbergman/docker-registry-cli/config"
	"github.com/pbergman/docker-registry-cli/helpers"
	"github.com/pbergman/docker-registry-cli/http"
	"github.com/pbergman/docker-registry-cli/logger"
)

// Mutex locked wrapper for the blob list
// so it can be safely used in goroutines
type lockedBlobList struct {
	list *blobList
	lock sync.Mutex
}

func NewLockedBlobList() *lockedBlobList {
	list := blobList(make(map[string][]string, 0))
	return &lockedBlobList{
		list: &list,
	}
}

func (l *lockedBlobList) get(blob string) []string {
	return (*l.list)[blob]
}

func (l *lockedBlobList) has(blob string) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.list.has(blob)
}

func (l *lockedBlobList) add(blob, repos, tag string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.list.add(blob, repos, tag)
}

// BlobList is a map with the digest as key and repository/tag as value
type blobList map[string][]string

func (b *blobList) has(blob string) bool {
	_, ok := (*b)[blob]
	return ok
}

func (b *blobList) add(blob, repos, tag string) {

	if false == b.has(blob) {
		(*b)[blob] = []string{repos + ":" + tag}
	} else {
		for _, name := range (*b)[blob] {
			if name == repos+":"+tag {
				return
			}
		}
		(*b)[blob] = append((*b)[blob], repos+":"+tag)
	}
}

func Delete(repository, tag string, dry bool) {

	if dry {
		fmt.Println("Running a dry-run.")
	}

	var wg sync.WaitGroup
	blobs := NewLockedBlobList()

	for _, repos := range *GetList() {
		for _, tagName := range repos.Tags {
			if repos.Name == repository && tagName == tag {
				continue
			}
			wg.Add(1)
			go func(list *lockedBlobList, repos, tag string) {
				manifests := GetManifest(repos, tag, true)
				for _, blob := range manifests.FsLayers {
					list.add(blob["blobSum"], repos, tag)
				}
				wg.Done()
			}(blobs, repos.Name, tagName)
		}
	}
	wg.Wait()

	manifest := GetManifest(repository, tag, true)
	remove := blobList(make(map[string][]string, 0))

	if true != dry {
		for _, layer := range manifest.FsLayers {
			if false == blobs.has(layer["blobSum"]) {
				logger.Logger.Debug(fmt.Sprintf("Adding layer %s to queued for removal", layer["blobSum"]))
				remove.add(layer["blobSum"], repository, tag)
			} else {
				logger.Logger.Debug(fmt.Sprintf("Skipping layer %s is used by: %s", layer["blobSum"], strings.Join(blobs.get(layer["blobSum"]), ", ")))
			}
		}
	} else {
		// For dry run we just print a overview
		table := helpers.NewTable("DIGEST", "REMOVING", "USED BY")
		for _, layer := range manifest.FsLayers {
			if false == blobs.has(layer["blobSum"]) {
				table.AddRow(layer["blobSum"], true)
			} else {
				table.AddRow(layer["blobSum"], false, strings.Join(blobs.get(layer["blobSum"]), ", "))
			}
		}
		table.Print()
	}

	if true != dry {

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

		first := true
		for layer, _ := range remove {
			wg.Add(1)
			// first one not so token can be cached
			if first {
				deleteBlob(layer)
				first = false
			} else {
				go deleteBlob(layer)
			}
		}
		wg.Wait()

		url := fmt.Sprintf("%s/v2/%s/manifests/%s", config.Config.RegistryHost, repository, manifest.Digest)
		logger.Logger.Debug("Requesting DELETE for url: " + url)
		req := http.NewRequest(url)
		req.SetMethod("DELETE")
		resp, err := req.Do()
		resp, err = resolve(resp, true)
		logger.Logger.CheckError(err)
	}
}
