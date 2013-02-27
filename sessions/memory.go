package sessions

import (
	"container/list"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/hujh/gmvc"
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
	watcher *memoryWatcher
}

func NewMemoryProvider(timeout time.Duration) *MemoryProvider {
	options := &Options{CookieName: defaultCookieName}
	storage := newMemoryStorage()
	watcher := newMemoryWatcher(storage, timeout)

	p := &MemoryProvider{
		Options: options,
		storage: storage,
		watcher: watcher,
	}

	runtime.SetFinalizer(p, func(x *MemoryProvider) { x.watcher.stop() })
	return p
}

func (p *MemoryProvider) GetSession(w http.ResponseWriter, r *http.Request, create bool) (session gmvc.Session, err error) {
	var store *memoryStore

	id := p.getCookieId(r)
	if id == "" {
		if create {
			id, err = p.generateId()
			if err != nil {
				return
			}
			for i := 0; i < 5; i++ {
				s, ok := p.storage.alloc(id)
				if ok {
					store = s
					break
				}
			}
			if store == nil {
				return nil, errors.New("fail to create session")
			}
			p.setCookieId(w, id, true)
		} else {
			return
		}
	} else {
		store = p.storage.touch(id)
	}

	session = &memorySession{
		provider: p,
		store:    store,
		id:       id,
		w:        w,
		valid:    true,
	}

	p.watcher.start()

	return
}

func (p *MemoryProvider) deleteSession(w http.ResponseWriter, id string) {
	p.storage.delete(id)
	p.setCookieId(w, id, false)
	return
}

func (p *MemoryProvider) generateId() (string, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		return "", err
	}

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

	return hex.EncodeToString(u), nil
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
	store    *memoryStore
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
	return s.store.set(key, value)
}

func (s *memorySession) Get(key string) (interface{}, error) {
	return s.store.get(key)
}

func (s *memorySession) Del(key string) error {
	return s.store.del(key)
}

type memoryStorage struct {
	mutex sync.Mutex
	elems map[string]*list.Element
	order *list.List
}

func newMemoryStorage() *memoryStorage {
	return &memoryStorage{
		elems: make(map[string]*list.Element),
		order: list.New(),
	}
}

func (m *memoryStorage) alloc(id string) (s *memoryStore, ok bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	e := m.elems[id]
	if e == nil {
		s = m.make(id)
		e := m.order.PushFront(s)
		m.elems[id] = e
		ok = true
		return
	}

	return nil, false
}

func (m *memoryStorage) touch(id string) (s *memoryStore) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	e := m.elems[id]
	if e == nil {
		s = m.make(id)
		e := m.order.PushFront(s)
		m.elems[id] = e
	} else {
		s = e.Value.(*memoryStore)
		s.utime = time.Now()
		m.order.MoveToFront(e)
	}

	return
}

func (m *memoryStorage) make(id string) *memoryStore {
	return &memoryStore{
		id:     id,
		utime:  time.Now(),
		values: make(map[string]interface{}),
	}
}

func (m *memoryStorage) delete(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	e := m.elems[id]
	if e != nil {
		delete(m.elems, id)
		m.order.Remove(e)
	}
}

func (m *memoryStorage) gc(timeout time.Duration) bool {
	performed := false
	for {
		if ok, count := m.cleanup(timeout, 1000); ok {
			performed = true
			if count == 0 {
				break
			}
		} else {
			break
		}
	}
	return performed
}

func (m *memoryStorage) cleanup(timeout time.Duration, limit int) (bool, int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.elems) == 0 {
		return false, 0
	}

	count := 0
	now := time.Now()
	for e := m.order.Back(); e != nil; {
		store := e.Value.(*memoryStore)
		if now.Sub(store.utime) > timeout {
			olde := e
			e = e.Prev()
			m.order.Remove(olde)
			delete(m.elems, store.id)
			count++
			if count >= limit {
				break
			}
		} else {
			break
		}
	}

	return true, count
}

type memoryStore struct {
	mutex  sync.RWMutex
	id     string
	utime  time.Time
	values map[string]interface{}
}

func (s *memoryStore) set(key string, value interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.values[key] = value
	return nil
}

func (s *memoryStore) get(key string) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.values[key], nil
}

func (s *memoryStore) del(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.values, key)
	return nil
}

type memoryWatcher struct {
	storage *memoryStorage
	timeout time.Duration
	sem     int32
	maxidle int
}

func newMemoryWatcher(storage *memoryStorage, timeout time.Duration) *memoryWatcher {
	if timeout == 0 {
		timeout = defaultTimeout
	}

	return &memoryWatcher{
		storage: storage,
		timeout: timeout,
		maxidle: 60,
	}
}

func (w *memoryWatcher) start() {
	if atomic.CompareAndSwapInt32(&w.sem, 0, 1) {
		go w.run()
	}
}

func (w *memoryWatcher) run() {
	ticker := time.NewTicker(time.Second)
	idle := 0
	for {
		if atomic.LoadInt32(&w.sem) == -1 {
			ticker.Stop()
			return
		}
		<-ticker.C
		if w.storage.gc(w.timeout) {
			idle = 0
		} else {
			idle++
			if idle >= w.maxidle {
				break
			}
		}
	}
	ticker.Stop()
	atomic.CompareAndSwapInt32(&w.sem, 1, 0)
}

func (w *memoryWatcher) stop() {
	atomic.StoreInt32(&w.sem, -1)
}
