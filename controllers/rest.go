package controllers

import (
	"fmt"
	"gmvc"
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
	n := t.NumMethod()

	if n == 0 {
		return fmt.Errorf("controller '%s' has no method", t)
	}

	router, err := r.router.Subrouter(pattern)
	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		m := t.Method(i)
		p := fmt.Sprintf("%s /", strings.ToUpper(m.Name))

		h, err := newMethodHandler(controller, m, r.arguments)
		if err != nil {
			return err
		}

		if err := router.Handle(p, h); err != nil {
			return err
		}
	}

	return nil
}
