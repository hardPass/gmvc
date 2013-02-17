package views

import (
	"fmt"
	"github.com/hujh/gmvc"
	"strings"
)

type Mux struct {
	views map[string]gmvc.View
}

func NewMux() *Mux {
	return &Mux{
		views: make(map[string]gmvc.View),
	}
}

func (m *Mux) Render(c *gmvc.Context, name string, value interface{}) error {
	var vname, subname string

	ns := strings.SplitN(name, ":", 2)

	if len(ns) > 1 {
		vname = ns[0]
		subname = ns[1]
	} else {
		vname = name
		subname = ""
	}

	view := m.views[vname]
	if view == nil {
		err := fmt.Errorf("not found view: %s", vname)
		return err
	}

	return view.Render(c, subname, value)
}

func (m *Mux) Set(name string, view gmvc.View) {
	m.views[name] = view
}

func (m *Mux) Del(name string) {
	delete(m.views, name)
}

func (m *Mux) Get(name string) gmvc.View {
	return m.views[name]
}
