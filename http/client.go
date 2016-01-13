package http

import "net/http"

var Client *http.Client

func init() {
	Client = &http.Client{}
}
