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
	defaultTimeout    time.Duration = 30 * time.Second
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
	id := p.readCookie(r)
	var store *memoryStore

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
			p.writeCookie(w, id)
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
	p.removeCookie(w, id)
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
		Name:     p.CookieName,
		Domain:   p.CookieDomain,
		Path:     p.CookiePath,
		Secure:   p.CookieSecure,
		HttpOnly: p.CookieHttpOnly,
		MaxAge:   p.CookieMaxAge,
		Value:    id,
	}

	w.Header().Add("Cache-Control", "private")
	http.SetCookie(w, cookie)
}

func (p *MemoryProvider) removeCookie(w http.ResponseWriter, id string) {
	cookie := &http.Cookie{
		Name:     p.CookieName,
		Domain:   p.CookieDomain,
		Path:     p.CookiePath,
		Secure:   p.CookieSecure,
		HttpOnly: p.CookieHttpOnly,
		MaxAge:   -1,
		Value:    id,
	}

	w.Header().Add("Cache-Control", "private")
	http.SetCookie(w, cookie)
}

type memorySession struct {
	provider *MemoryProvider
	store    *memoryStore
	id       string
	valid    bool
	w        http.ResponseWriter
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

func (m *memoryStorage) cleanup(timeout time.Duration) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.elems) == 0 {
		return false
	}

	now := time.Now()
	for e := m.order.Back(); e != nil; e = e.Prev() {
		store := e.Value.(*memoryStore)
		if store.expire(now, timeout) {
			delete(m.elems, store.id)
			olde := e
			e = e.Prev()
			m.order.Remove(olde)
		}
		break
	}

	return true
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

func (s *memoryStore) expire(t time.Time, d time.Duration) bool {
	return t.Sub(s.utime) > d
}

type memoryWatcher struct {
	storage *memoryStorage
	timeout time.Duration
	sem     int32
	maxidle int
	ticker  *time.Ticker
}

func newMemoryWatcher(storage *memoryStorage, timeout time.Duration) *memoryWatcher {
	if timeout == 0 {
		timeout = defaultTimeout
	}

	return &memoryWatcher{
		storage: storage,
		timeout: timeout,
		maxidle: 10,
	}
}

func (w *memoryWatcher) start() {
	if atomic.CompareAndSwapInt32(&w.sem, 0, 1) {
		go w.run()
	}
}

func (w *memoryWatcher) run() {
	w.ticker = time.NewTicker(time.Second)
	idle := 0
	for {
		<-w.ticker.C
		if w.storage.cleanup(w.timeout) {
			idle = 0
		} else {
			idle++
			if idle >= w.maxidle {
				break
			}
		}
	}
	w.ticker.Stop()
	atomic.CompareAndSwapInt32(&w.sem, 1, 0)
}

func (w *memoryWatcher) stop() {
	atomic.StoreInt32(&w.sem, -1)
	if w.ticker != nil {
		w.ticker.Stop()
	}
}
