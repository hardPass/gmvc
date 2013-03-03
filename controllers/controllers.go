package controllers

import (
	"fmt"
	"github.com/hujh/gmvc"
	"reflect"
	"regexp"
	"strings"
)

var (
	regexMapping  = regexp.MustCompile("^((?:\\S+\\s+){0,1}/\\S*)\\s+(\\S+)$")
	mappingSyntax = "[HttpMethods] <UrlPattern> <ControllerMethod>"
)

type Controller interface {
	RequestMapping() string
}

func Register(router *gmvc.Router, pattern string, controller Controller) error {
	router, err := router.Subrouter(pattern)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(controller)
	mapping := controller.RequestMapping()
	for _, line := range strings.Split(mapping, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		match := regexMapping.FindStringSubmatch(line)
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
