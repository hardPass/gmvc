package controllers

import (
	"fmt"
	"github.com/hujh/gmvc"
	"reflect"
	"strings"
)

type REST struct {
	router    *gmvc.Router
	arguments []Argument
}

func NewREST(router *gmvc.Router) *REST {
	r := &REST{
		router:    router,
		arguments: make([]Argument, len(defaultArguments)),
	}
	copy(r.arguments, defaultArguments)
	return r
}

func (r *REST) Register(pattern string, controller interface{}) error {
	t := reflect.TypeOf(controller)
	nm := t.NumMethod()

	if nm == 0 {
		return fmt.Errorf("controller '%s' has no method", t)
	}

	router, err := r.router.Subrouter(pattern)
	if err != nil {
		return err
	}

	reg := false
	for i := 0; i < nm; i++ {
		m := t.Method(i)
		if strings.Title(m.Name) != m.Name {
			continue
		}

		p := fmt.Sprintf("%s /", strings.ToUpper(m.Name))
		h, err := newMethodHandler(controller, m, r.arguments)
		if err != nil {
			return err
		}

		if err := router.Handle(p, h); err != nil {
			return err
		}
		reg = true
	}

	if !reg {
		return fmt.Errorf("controller '%s' has no accessible method", t)
	}

	return nil
}
