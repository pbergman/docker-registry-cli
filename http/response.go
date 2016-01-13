package http

import (
	"errors"
	"net/http"
	"reflect"
	"regexp"
)

type Response struct {
	*http.Response
}

func (r *Response) GetAuthChallenge() (*AuthChallenge, error) {
	challenge := &AuthChallenge{}
	header := r.Header.Get("Www-Authenticate")
	if len(header) <= 0 {
		return nil, errors.New("No or empty 'Www-Authenticate' header")
	}
	regex := regexp.MustCompile(`([^\s]+)\s|([^=]+)="([^"]+)",?`)
	match := regex.FindAllStringSubmatchIndex(header, -1)
	if len(match) <= 0 {
		return nil, errors.New("Invalid 'Www-Authenticate' header")
	}
	challenge.Scheme = header[match[0:1][0][2]:match[0:1][0][3]]
	val := reflect.ValueOf(challenge).Elem()
	for i := 0; i < val.NumField(); i++ {
		for c := 1; c < len(match); c++ {
			if header[match[c : c+1][0][4]:match[c : c+1][0][5]] == val.Type().Field(i).Tag.Get("header") {
				val.Field(i).SetString(header[match[c : c+1][0][6]:match[c : c+1][0][7]])
			}
		}
	}
	if challenge.Error != "" {
		return nil, errors.New(challenge.Error)
	}
	return challenge, nil
}
