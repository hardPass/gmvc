package gmvc

import (
	"errors"
	"net/http"
)

type Session interface {
	Id() string

	Valid() bool
	Invalidate() error

	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Del(key string) error
}

type SessionStorage interface {
	Create() (string, error)
	Remove(id string) error
	Touch(id string, create bool) (bool, error)

	SetAttr(id string, key string, value interface{}) error
	GetAttr(id string, key string) (interface{}, error)
	DelAttr(id string, key string) error
}

type SessionManager struct {
	Storage        SessionStorage
	CookieName     string
	CookiePath     string
	CookieDomain   string
	CookieMaxAge   int
	CookieSecure   bool
	CookieHttpOnly bool
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		Storage:    nil,
		CookieName: "sid",
	}
}

func (sm *SessionManager) Get(w http.ResponseWriter, r *http.Request, create bool) (s Session, err error) {
	if sm == nil {
		if create {
			return nil, errors.New("nil session manager")
		}
		return nil, nil
	}

	id := sm.readCookie(r)
	if id == "" {
		if create {
			id, err = sm.Storage.Create()
			if err != nil {
				return nil, err
			}
			if id == "" {
				return nil, errors.New("generate empty sessionid")
			}
			sm.writeCookie(w, id)
		} else {
			return nil, nil
		}
	} else {
		exists, err := sm.Storage.Touch(id, true)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, nil
		}
	}

	return newSession(sm, id, w), nil
}

func (sm *SessionManager) Del(w http.ResponseWriter, id string) error {
	if err := sm.Storage.Remove(id); err != nil {
		return err
	}

	sm.removeCookie(w, id)
	return nil
}

func (sm *SessionManager) readCookie(r *http.Request) string {
	for _, c := range r.Cookies() {
		if c.Name == sm.CookieName {
			return c.Value
		}
	}

	return ""
}

func (sm *SessionManager) writeCookie(w http.ResponseWriter, id string) {
	cookie := &http.Cookie{
		Value:    id,
		Name:     sm.CookieName,
		Domain:   sm.CookieDomain,
		Path:     sm.CookiePath,
		Secure:   sm.CookieSecure,
		HttpOnly: sm.CookieHttpOnly,
		MaxAge:   sm.CookieMaxAge,
	}

	http.SetCookie(w, cookie)
	w.Header().Add("Cache-Control", "private")
}

func (sm *SessionManager) removeCookie(w http.ResponseWriter, id string) {
	cookie := &http.Cookie{
		Value:    id,
		Name:     sm.CookieName,
		Domain:   sm.CookieDomain,
		Path:     sm.CookiePath,
		Secure:   sm.CookieSecure,
		HttpOnly: sm.CookieHttpOnly,
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)
	w.Header().Add("Cache-Control", "private")
}

type session struct {
	id string
	m  *SessionManager
	w  http.ResponseWriter
	v  bool
}

func newSession(m *SessionManager, id string, w http.ResponseWriter) *session {
	return &session{
		id: id,
		m:  m,
		w:  w,
		v:  true,
	}
}

func (s *session) Id() string {
	return s.id
}

func (s *session) Valid() bool {
	return s.v
}

func (s *session) Invalidate() error {
	s.v = false
	return s.m.Del(s.w, s.id)
}

func (s *session) Set(key string, value interface{}) error {
	return s.m.Storage.SetAttr(s.id, key, value)
}

func (s *session) Get(key string) (interface{}, error) {
	return s.m.Storage.GetAttr(s.id, key)
}

func (s *session) Del(key string) error {
	return s.m.Storage.DelAttr(s.id, key)
}
