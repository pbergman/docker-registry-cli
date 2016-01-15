package http

import (
	"net/http"
)

type Request struct {
	request *http.Request
}

func FromRequest(r *http.Request) *Request {
	return &Request{request: r}
}

func NewRequest(url string) *Request {
	request, _ := http.NewRequest("GET", url, nil)
	return &Request{request: request}
}

func (r *Request) Do() (*Response, error) {

	resp, err := Client.Do(r.request)

	if err != nil {
		return nil, err
	} else {
		return &Response{resp}, nil
	}
}

func (r *Request) Raw() *http.Request {
	return r.request
}

func (r *Request) SetMethod(method string) {
	r.request.Method = method
}

func (r *Request) AddHeader(name, value string) {
	r.request.Header.Add(name, value)
}

func (r *Request) SetBasicAuth(username, password string) {
	r.request.SetBasicAuth(username, password)
}

func (r *Request) AddBearerToken(token string) {
	r.request.Header.Set("Authorization", "Bearer " + token)
}
