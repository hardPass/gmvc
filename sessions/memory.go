package sessions

import (
	"container/list"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/hujh/gmvc"
	"io"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	defaultCookieName string        = "gsessionid"
	defaultTimeout    time.Duration = 30 * time.Minute
)

type Options struct {
	CookieName     string
	CookiePath     string
	CookieDomain   string
	CookieMaxAge   int
	CookieSecure   bool
	CookieHttpOnly bool
}

type MemoryProvider struct {
	*Options
	storage *memoryStorage
}

func NewMemoryProvider(timeout time.Duration) *MemoryProvider {
	options := &Options{CookieName: defaultCookieName, CookieHttpOnly: true}
	storage := newMemoryStorage(timeout)

	p := &MemoryProvider{
		Options: options,
		storage: storage,
	}

	runtime.SetFinalizer(p, func(x *MemoryProvider) { x.storage.destory() })
	return p
}

func (p *MemoryProvider) GetSession(w http.ResponseWriter, r *http.Request, create bool) (s gmvc.Session, err error) {
	var values *memoryValues

	id := p.getCookieId(r)
	if id == "" {
		if create {
			for i := 0; i < 5; i++ {
				id, err = p.generateId()
				if err != nil {
					return
				}
				v, ok := p.storage.alloc(id)
				if ok {
					values = v
					break
				}
			}
			if values == nil {
				return nil, errors.New("fail to create session")
			}
			p.setCookieId(w, id, true)
		} else {
			return
		}
	} else {
		values = p.storage.touch(id)
	}

	s = &memorySession{
		provider: p,
		values:   values,
		id:       id,
		w:        w,
		valid:    true,
	}

	p.storage.gc()

	return
}

func (p *MemoryProvider) deleteSession(w http.ResponseWriter, id string) {
	p.storage.delete(id)
	p.setCookieId(w, id, false)
	return
}

func (p *MemoryProvider) generateId() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}

	hash := md5.New()
	if _, err := hash.Write(b); err != nil {
		return "", err
	}

	b = hash.Sum(nil)
	return hex.EncodeToString(b), nil
}

func (p *MemoryProvider) getCookieId(r *http.Request) string {
	for _, c := range r.Cookies() {
		if c.Name == p.CookieName {
			return c.Value
		}
	}
	return ""
}

func (p *MemoryProvider) setCookieId(w http.ResponseWriter, id string, persist bool) {
	cookie := &http.Cookie{
		Name:     p.CookieName,
		Domain:   p.CookieDomain,
		Path:     p.CookiePath,
		Secure:   p.CookieSecure,
		HttpOnly: p.CookieHttpOnly,
		Value:    id,
	}

	if persist {
		cookie.MaxAge = p.CookieMaxAge
	} else {
		cookie.MaxAge = -1
	}

	w.Header().Add("Cache-Control", "private")
	http.SetCookie(w, cookie)
}

type memorySession struct {
	w        http.ResponseWriter
	provider *MemoryProvider
	values   *memoryValues
	id       string
	valid    bool
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

func (s *memorySession) Set(key string, value interface{}) error {
	return s.values.set(key, value)
}

func (s *memorySession) Get(key string) (interface{}, error) {
	return s.values.get(key)
}

func (s *memorySession) Del(key string) error {
	return s.values.del(key)
}

type memoryStorage struct {
	mutex   sync.Mutex
	elems   map[string]*list.Element
	order   *list.List
	timeout time.Duration
	maxidle int
	sem     int32
}

func newMemoryStorage(timeout time.Duration) *memoryStorage {
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	return &memoryStorage{
		elems:   make(map[string]*list.Element),
		order:   list.New(),
		timeout: timeout,
		maxidle: 60,
	}
}

func (s *memoryStorage) alloc(id string) (v *memoryValues, ok bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	e := s.elems[id]
	if e == nil {
		v = s.make(id)
		e := s.order.PushFront(v)
		s.elems[id] = e
		ok = true
		return
	}

	return nil, false
}

func (s *memoryStorage) touch(id string) (v *memoryValues) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	e := s.elems[id]
	if e == nil {
		v = s.make(id)
		e := s.order.PushFront(v)
		s.elems[id] = e
	} else {
		v = e.Value.(*memoryValues)
		v.utime = time.Now()
		s.order.MoveToFront(e)
	}

	return
}

func (s *memoryStorage) make(id string) *memoryValues {
	return &memoryValues{
		id:     id,
		utime:  time.Now(),
		values: make(map[string]interface{}),
	}
}

func (s *memoryStorage) delete(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if e := s.elems[id]; e != nil {
		delete(s.elems, id)
		s.order.Remove(e)
	}
}

func (s *memoryStorage) gc() {
	if atomic.CompareAndSwapInt32(&s.sem, 0, 1) {
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			idle := 0
			timeout := s.timeout
			for {
				if atomic.LoadInt32(&s.sem) == -1 {
					return
				}

				<-ticker.C
				performed := false

				for {
					if ok, count := s.trim(timeout, 4096); ok {
						performed = true
						if count == 0 {
							break
						}
					} else {
						break
					}
				}

				if performed {
					idle = 0
				} else {
					idle++
					if idle >= s.maxidle {
						break
					}
				}
			}
			atomic.CompareAndSwapInt32(&s.sem, 1, 0)
		}()
	}
}

func (s *memoryStorage) trim(timeout time.Duration, limit int) (bool, int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.elems) == 0 {
		return false, 0
	}

	count := 0
	now := time.Now()
	for e := s.order.Back(); e != nil; {
		values := e.Value.(*memoryValues)
		if now.Sub(values.utime) > timeout {
			if count > limit {
				break
			}
			olde := e
			e = e.Prev()
			s.order.Remove(olde)
			delete(s.elems, values.id)
			count++

		} else {
			break
		}
	}

	return true, count
}

func (s *memoryStorage) destory() {
	atomic.StoreInt32(&s.sem, -1)
}

type memoryValues struct {
	mutex  sync.RWMutex
	id     string
	utime  time.Time
	values map[string]interface{}
}

func (s *memoryValues) set(key string, value interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.values[key] = value
	return nil
}

func (s *memoryValues) get(key string) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.values[key], nil
}

func (s *memoryValues) del(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.values, key)
	return nil
}
