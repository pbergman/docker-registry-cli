package account

type User struct {
	Username string `json:"name,omitempty"`
	Password string `json:"pass,omitempty "`
}

func NewEmptyUser() *User {
	return &User{}
}

func NewUser(username, password string) *User {
	return &User{username, password}
}
