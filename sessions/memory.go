package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	defaultTimeout time.Duration = 30 * time.Second
)

type entry struct {
	values map[string]interface{}
	utime  time.Time
}

type MemoryStorage struct {
	mutex    sync.RWMutex
	sessions map[string]*entry
	timeout  time.Duration
	wsem     int32
	wticker  *time.Ticker
	wmaxidle int
}

func NewMemoryStorage(timeout time.Duration) *MemoryStorage {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return &MemoryStorage{
		sessions: make(map[string]*entry),
		timeout:  timeout,
		wsem:     0,
		wmaxidle: 10,
	}
}

func (s *MemoryStorage) Create() (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := ""

	u := make([]byte, 16)
	for i := 0; i < 5; i++ {
		_, err := rand.Read(u)
		if err != nil {
			return "", err
		}

		u[8] = (u[8] | 0x80) & 0xBF
		u[6] = (u[6] | 0x40) & 0x4F

		t := hex.EncodeToString(u)

		if _, found := s.sessions[id]; !found {
			id = t
			break
		}
	}

	if id == "" {
		return "", errors.New("fail to generate sessionid")
	}

	entry := &entry{
		values: make(map[string]interface{}),
		utime:  time.Now(),
	}

	s.sessions[id] = entry
	s.watch()

	return id, nil
}

func (s *MemoryStorage) Remove(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.sessions, id)
	return nil
}

func (s *MemoryStorage) Touch(id string, create bool) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	se := s.sessions[id]

	if se == nil {
		if create {
			se = &entry{
				values: make(map[string]interface{}),
				utime:  time.Now(),
			}
			s.sessions[id] = se
			s.watch()
			return true, nil
		}
		return false, nil
	}

	se.utime = time.Now()
	return true, nil
}

func (s *MemoryStorage) SetAttr(id string, key string, value interface{}) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	entry := s.sessions[id]

	if entry == nil {
		return errors.New("not found session: " + id)
	}

	entry.values[key] = value
	return nil
}

func (s *MemoryStorage) GetAttr(id string, key string) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	entry := s.sessions[id]

	if entry == nil {
		return nil, errors.New("not found session: " + id)
	}

	return entry.values[key], nil
}

func (s *MemoryStorage) DelAttr(id string, key string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	entry := s.sessions[id]
	if entry == nil {
		return errors.New("not found session: " + id)
	}

	delete(entry.values, key)
	return nil
}

func (s *MemoryStorage) watch() {
	if atomic.CompareAndSwapInt32(&s.wsem, 0, 1) {
		go func() {
			s.wticker = time.NewTicker(time.Second)
			idle := 0
			for {
				<-s.wticker.C
				now := time.Now()
				s.mutex.Lock()
				if len(s.sessions) == 0 {
					idle++
					if idle >= s.wmaxidle {
						s.mutex.Unlock()
						break
					}
				} else {
					idle = 0
					for id, entry := range s.sessions {
						if now.Sub(entry.utime) > s.timeout {
							delete(s.sessions, id)
						}
					}
				}
				s.mutex.Unlock()
			}
			s.wticker.Stop()
			atomic.StoreInt32(&s.wsem, 0)
		}()
	}
}
