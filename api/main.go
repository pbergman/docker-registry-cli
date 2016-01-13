package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pbergman/docker-registery-cli/config"
	"github.com/pbergman/docker-registery-cli/http"
	"github.com/pbergman/docker-registery-cli/logger"
	base_logger "github.com/pbergman/logger"
)

const (
	// https://docs.docker.com/registry/spec/api/#errors-2
	BLOB_UNKNOWN          string = "BLOB_UNKNOWN"
	BLOB_UPLOAD_INVALID   string = "BLOB_UPLOAD_INVALID"
	BLOB_UPLOAD_UNKNOWN   string = "BLOB_UPLOAD_UNKNOWN"
	DIGEST_INVALID        string = "DIGEST_INVALID"
	MANIFEST_BLOB_UNKNOWN string = "MANIFEST_BLOB_UNKNOWN"
	MANIFEST_INVALID      string = "MANIFEST_INVALID"
	MANIFEST_UNKNOWN      string = "MANIFEST_UNKNOWN"
	MANIFEST_UNVERIFIED   string = "MANIFEST_UNVERIFIED"
	NAME_INVALID          string = "NAME_INVALID"
	NAME_UNKNOWN          string = "NAME_UNKNOWN"
	SIZE_INVALID          string = "SIZE_INVALID"
	TAG_INVALID           string = "TAG_INVALID"
	UNAUTHORIZED          string = "UNAUTHORIZED"
	DENIED                string = "DENIED"
	UNSUPPORTED           string = "UNSUPPORTED"
)

func resolve(response *http.Response, exit_on_failure bool) (*http.Response, error) {

	logger.Logger.Debug("Status: " + response.Status)

	switch response.StatusCode {
	case 200, 202:
		return response, nil
	case 401:
		// Make sure we got a user
		config.Config.CheckUser()
		// Get Authenticate chalange
		challenge, err := response.GetAuthChallenge()
		logger.Logger.CheckError(err)
		// Will fetch token and authenticate user
		token, err := config.Config.TokenManager.GetToken(challenge)
		logger.Logger.CheckError(err)
		// Build new request and add bearer token
		request := http.FromRequest(response.Request)
		request.AddBearerToken(token.Token)
		resp, err := request.Do()
		logger.Logger.CheckError(err)
		// recheck response
		return resolve(resp, exit_on_failure)
	default:
		err := GetErrorResponse(response)
		if exit_on_failure {
			DisplayErrorResponse(err)
			os.Exit(1)
		}
		return nil, err
	}
}

type ErrorsResponse struct {
	Errors     []*Error `json:"errors"`
	StatusCode int
}

func (e *ErrorsResponse) Error() string {
	return fmt.Sprintf("[%s] %s", e.Errors[0].Code, e.Errors[0].Message)
}

type Error struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail,omitempty"`
}

func DisplayErrorResponse(err *ErrorsResponse) {
	logger.Logger.Error(base_logger.NewContextMessage(err.Error(), err.Errors[0].Detail.(map[string]interface{})))
}

func GetErrorResponse(response *http.Response) *ErrorsResponse {
	body, err := ioutil.ReadAll(response.Body)
	logger.Logger.CheckError(err)
	errResponse := &ErrorsResponse{}

	if strings.HasPrefix(response.Header.Get("Content-Type"), "application/json") {
		err = json.Unmarshal(body, errResponse)
		logger.Logger.CheckWarning(err)
	} else {
		errResponse.Errors = make([]*Error, 1)
		errResponse.Errors[0] = &Error{
			Message: strings.TrimSpace(string(body)),
			Code:    "UNKOWN",
			Detail:  make(map[string]interface{}, 0),
		}
	}
	// No JSON message?
	if err != nil {
		errResponse.Errors = make([]*Error, 1)
		errResponse.Errors[0] = &Error{
			Message: strings.TrimSpace(string(body)),
			Code:    "UNKOWN",
			Detail:  make(map[string]interface{}, 0),
		}
	}
	errResponse.StatusCode = response.StatusCode
	return errResponse
}
