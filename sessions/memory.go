package sessions

import (
	"container/list"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/hujh/gmvc"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	defaultCookieName string        = "gsessionid"
	defaultTimeout    time.Duration = 30 * time.Second
)

type CookieOptions struct {
	CookieName     string
	CookiePath     string
	CookieDomain   string
	CookieMaxAge   int
	CookieSecure   bool
	CookieHttpOnly bool
}

type MemoryProvider struct {
	*CookieOptions

	mutex    sync.RWMutex
	storage  map[string]*list.Element
	order    *list.List
	timeout  time.Duration
	wsem     int32
	wticker  *time.Ticker
	wmaxidle int
}

func NewMemoryProvider(timeout time.Duration) *MemoryProvider {
	if timeout == 0 {
		timeout = defaultTimeout
	}

	option := &CookieOptions{
		CookieName: defaultCookieName,
	}

	return &MemoryProvider{
		CookieOptions: option,
		storage:       make(map[string]*list.Element),
		order:         list.New(),
		timeout:       timeout,
		wsem:          0,
		wmaxidle:      10,
	}
}

func (p *MemoryProvider) GetSession(w http.ResponseWriter, r *http.Request, create bool) (gmvc.Session, error) {
	id := p.readCookie(r)

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if id == "" {
		if create {
			id, err := p.generateId()
			if err != nil {
				return nil, err
			}
			p.writeCookie(w, id)
		} else {
			return nil, nil
		}
	}

	store := p.store(id)
	session := &memorySession{
		id:          id,
		w:           w,
		valid:       true,
		memoryStore: store,
	}

	return session, nil
}

func (p *MemoryProvider) deleteSession(w http.ResponseWriter, id string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	delete(p.storage, id)
	p.removeCookie(w, id)
	return
}

func (p *MemoryProvider) generateId() (string, error) {
	id := ""
	u := make([]byte, 16)
	for i := 0; i < 5; i++ {
		_, err := rand.Read(u)
		if err != nil {
			return "", err
		}
		u[8] = (u[8] | 0x80) & 0xBF
		u[6] = (u[6] | 0x40) & 0x4F
		str := hex.EncodeToString(u)
		if _, ok := p.storage[id]; !ok {
			id = str
			break
		}
	}

	if id == "" {
		return "", errors.New("fail to generate sessionid")
	}

	return id, nil
}

func (p *MemoryProvider) store(id string) (s *memoryStore) {
	e := p.storage[id]

	if e == nil {
		s = &memoryStore{
			id:     id,
			utime:  time.Now(),
			values: make(map[string]interface{}),
		}
		e = p.order.PushFront(s)
		p.storage[id] = e
	} else {
		s = e.Value.(*memoryStore)
		s.utime = time.Now()
		p.order.MoveToFront(e)
	}

	return s
}

func (p *MemoryProvider) readCookie(r *http.Request) string {
	for _, c := range r.Cookies() {
		if c.Name == p.CookieName {
			return c.Value
		}
	}

	return ""
}

func (p *MemoryProvider) writeCookie(w http.ResponseWriter, id string) {
	cookie := &http.Cookie{
		Value:    id,
		Name:     p.CookieName,
		Domain:   p.CookieDomain,
		Path:     p.CookiePath,
		Secure:   p.CookieSecure,
		HttpOnly: p.CookieHttpOnly,
		MaxAge:   p.CookieMaxAge,
	}

	http.SetCookie(w, cookie)
	w.Header().Add("Cache-Control", "private")
}

func (p *MemoryProvider) removeCookie(w http.ResponseWriter, id string) {
	cookie := &http.Cookie{
		Value:    id,
		Name:     p.CookieName,
		Domain:   p.CookieDomain,
		Path:     p.CookiePath,
		Secure:   p.CookieSecure,
		HttpOnly: p.CookieHttpOnly,
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)
	w.Header().Add("Cache-Control", "private")
}

func (p *MemoryProvider) watch() {
	if atomic.CompareAndSwapInt32(&p.wsem, 0, 1) {
		go p.cleanup()
	}
}

func (p *MemoryProvider) cleanup() {
	p.wticker = time.NewTicker(time.Second)
	idle := 0
	for {
		<-p.wticker.C
		now := time.Now()
		p.mutex.Lock()
		if len(p.storage) == 0 {
			idle++
			if idle >= p.wmaxidle {
				p.mutex.Unlock()
				break
			}
		} else {
			idle = 0
			for e := p.order.Back(); e != nil; e = e.Prev() {
				store := e.Value.(*memoryStore)
				if store.expire(now, p.timeout) {
					delete(p.storage, store.id)
					oe := e
					e = e.Prev()
					p.order.Remove(oe)
				}
				break
			}
		}
		p.mutex.Unlock()
	}
	p.wticker.Stop()
	atomic.StoreInt32(&p.wsem, 0)
}

type memoryStore struct {
	mutex  sync.RWMutex
	id     string
	utime  time.Time
	values map[string]interface{}
}

func (s *memoryStore) Set(key string, value interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.values[key] = value
	return nil
}

func (s *memoryStore) Get(key string) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.values[key], nil
}

func (s *memoryStore) Del(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.values, key)
	return nil
}

func (s *memoryStore) expire(t time.Time, d time.Duration) bool {
	return t.Sub(s.utime) > d
}

type memorySession struct {
	provider *MemoryProvider
	id       string
	valid    bool
	w        http.ResponseWriter
	*memoryStore
}

func (s *memorySession) Id() string {
	return s.id
}

func (s *memorySession) Valid() bool {
	return s.valid
}

func (s *memorySession) Invalidate() error {
	if s.valid {
		s.valid = false
		s.provider.deleteSession(s.w, s.id)
	}
	return nil
}

func (s *memorySession) Save() error {
	return nil
}
