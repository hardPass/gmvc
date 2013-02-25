package gmvc

import (
	"net/http"
)

type SessionProvider interface {
	GetSession(w http.ResponseWriter, r *http.Request, create bool) (Session, error)
}

type Session interface {
	Id() string

	Valid() bool
	Invalidate() error

	Save() error

	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Del(key string) error
}
