package token

import "strings"

type Scope struct {
	Type   string
	Name   string
	Action []string
}

func NewScope(st string, sn string, a ...string) *Scope {
	return &Scope{
		Type:   st,
		Name:   sn,
		Action: a,
	}
}

func (s *Scope) AddAction(action string) {
	s.Action = append(s.Action, action)
}

func (s *Scope) HasActions(actions []string) bool {

	inAction := func(a string) bool {
		for _, action := range actions {
			if action == a {
				return true
			}
		}
		return false
	}

	for _, action := range s.Action {
		if false == inAction(action) {
			return false
		}
	}

	return true
}

func (s *Scope) HasAction(action string) bool {
	for _, a := range s.Action {
		if a == action {
			return true
		}
	}
	return false
}

func (s *Scope) String() string {
	return s.Type + ":" + s.Name + ":" + strings.Join(s.Action, ",")
}
