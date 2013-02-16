package views

import (
	"fmt"
	"gmvc"
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
	var subname, vname string

	s := strings.SplitN(name, ":", 2)

	if len(s) > 1 {
		subname = s[0]
		vname = s[1]
	} else {
		subname = name
		vname = ""
	}

	view := m.views[subname]
	if view == nil {
		err := fmt.Errorf("not found view: %s", subname)
		return err
	}

	return view.Render(c, vname, value)
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
