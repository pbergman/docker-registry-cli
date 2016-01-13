package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/pbergman/docker-registery-cli/config"
	"github.com/pbergman/docker-registery-cli/http"
	"github.com/pbergman/docker-registery-cli/logger"
)

type Header struct {
	Jwk map[string]string `json:"jwk"`
	Alg string            `json:"alg"`
}

type Signatures struct {
	Signature string  `json:"signature"`
	Protected string  `json:"protected"`
	header    *Header `json:"header"`
}

type History struct {
	V1Compatibility string `json:"v1Compatibility"`
}

func (h *History) Unpack() map[string]interface{} {
	s, err := strconv.Unquote("`" + h.V1Compatibility + "`")
	logger.Logger.CheckError(err)
	var unpacked map[string]interface{}
	err = json.Unmarshal([]byte(s), &unpacked)
	logger.Logger.CheckError(err)
	return unpacked
}

func (h *History) Print() *bytes.Buffer {
	s, err := strconv.Unquote("`" + h.V1Compatibility + "`")
	logger.Logger.CheckError(err)
	var buff bytes.Buffer
	err = json.Indent(&buff, []byte(s), "", "\t")
	logger.Logger.CheckError(err)
	return &buff
}

type Manifest struct {
	Digest        string
	SchemaVersion int                 `json:"schemaVersion"`
	Name          string              `json:"name"`
	Tag           string              `json:"tag"`
	Architecture  string              `json:"architecture"`
	FsLayers      []map[string]string `json:"fsLayers"`
	History       []*History          `json:"history"`
	Signatures    []*Signatures       `json:"signatures"`
}

func NewManifest() *Manifest {
	signatures := make([]*Signatures, 0)
	signatures = append(signatures, &Signatures{header: &Header{}})
	maifest := &Manifest{Signatures: signatures, History: make([]*History, 0)}
	return maifest
}

var manifests map[string]map[string]*Manifest

func init() {
	manifests = make(map[string]map[string]*Manifest, 0)
}

func GetManifest(repository, tag string) *Manifest {
	if manifest, ok := manifests[repository][tag]; ok {
		return manifest
	}
	url := fmt.Sprintf("%s/v2/%s/manifests/%s", config.Config.RegistryHost, repository, tag)
	logger.Logger.Debug("Requesting url: " + url)
	req := http.NewRequest(url)
	resp, err := req.Do()
	logger.Logger.CheckError(err)
	resp, err = resolve(resp, true)
	body, err := ioutil.ReadAll(resp.Body)
	logger.Logger.CheckError(err)
	manifest := NewManifest()
	err = json.Unmarshal(body, manifest)
	logger.Logger.CheckError(err)
	manifest.Digest = resp.Header.Get("Docker-Content-Digest")
	if _, ok := manifests[repository]; !ok {
		manifests[repository] = make(map[string]*Manifest, 0)
	}
	manifests[repository][tag] = manifest
	return manifest
}
