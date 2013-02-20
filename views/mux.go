package views

import (
	"fmt"
	"github.com/hujh/gmvc"
	"strings"
	"sync"
)

type Mux struct {
	mutex sync.RWMutex
	views map[string]gmvc.View
}

func NewMux() *Mux {
	return &Mux{
		views: make(map[string]gmvc.View),
	}
}

func (m *Mux) Render(c *gmvc.Context, name string, data interface{}) error {
	ns := strings.SplitN(name, ":", 2)
	var k, n string
	if len(ns) > 1 {
		k = ns[0]
		n = ns[1]
	} else {
		k = name
		n = ""
	}

	m.mutex.RLock()
	view := m.views[k]
	m.mutex.RUnlock()

	if view == nil {
		err := fmt.Errorf("view not found '%s'", k)
		return err
	}

	return view.Render(c, n, data)
}

func (m *Mux) Set(name string, view gmvc.View) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.views[name] = view
}

func (m *Mux) Del(name string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.views, name)
}

func (m *Mux) Get(name string) gmvc.View {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.views[name]
}
