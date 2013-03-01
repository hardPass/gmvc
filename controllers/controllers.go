package controllers

import (
	"fmt"
	"github.com/hujh/gmvc"
	"reflect"
	"regexp"
	"strings"
)

var (
	mapping       = regexp.MustCompile("^((?:\\S+\\s+){0,1}/\\S*)\\s+(\\S+)$")
	mappingSyntax = "[HttpMethods] <UrlPattern> <ControllerMethod>"
)

type Controller interface {
	RequestMapping() string
}

type Controllers struct {
	router *gmvc.Router
}

func New(router *gmvc.Router) *Controllers {
	return &Controllers{
		router: router,
	}
}

func (c *Controllers) Register(pattern string, controller Controller) error {
	t := reflect.TypeOf(controller)

	router, err := c.router.Subrouter(pattern)
	if err != nil {
		return err
	}

	lines := strings.Split(controller.RequestMapping(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		match := mapping.FindStringSubmatch(line)
		if match == nil {
			return fmt.Errorf("controller %s has incorrect format mapping: '%s', syntax: %s", t, line, mappingSyntax)
		}

		pattern := match[1]
		methodName := match[2]

		method, ok := t.MethodByName(methodName)
		if !ok {
			return fmt.Errorf("controller %s has incorrect mapping: '%s', reason: no method %s", t, line, methodName)
		}

		if strings.Title(method.Name) != method.Name {
			return fmt.Errorf("controller %s has incorrect mapping: '%s', reason: %s is not accessable", t, line, methodName)
		}

		handler, err := newMethodHandler(controller, method)
		if err != nil {
			return err
		}

		if err := router.Handle(pattern, handler); err != nil {
			return err
		}
	}

	return nil
}
