package token

import (
	"time"
)

type Token struct {
	Token   string    `json:"token"`
	Expires int       `json:"expires_in"`
	Issued  time.Time `json:"issued_at"`
}

func (t *Token) ExpireTime() time.Time {
	return t.Issued.Add(time.Second * time.Duration(t.Expires))
}

func (t *Token) IsValid() bool {
	return time.Now().Before(t.Issued.Add(time.Second * time.Duration(t.Expires)))
}
