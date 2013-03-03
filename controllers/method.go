package controllers

import (
	"github.com/hujh/gmvc"
	"net/http"
	"reflect"
)

var (
	typeOfError = reflect.TypeOf(new(error)).Elem()
)

type methodHandler struct {
	controller reflect.Value
	method     reflect.Method
	args       []argument
	outErr     int
}

func newMethodHandler(controller interface{}, method reflect.Method) (*methodHandler, error) {
	numIn := method.Type.NumIn()
	args := make([]argument, numIn-1)

	for i := 1; i < numIn; i++ {
		var arg argument
		t := method.Type.In(i)

		for _, a := range arguments {
			at := a.Type()
			if (at.Kind() == reflect.Interface && t.Implements(at)) || t == at {
				arg = a
				break
			}
		}

		if arg == nil {
			arg = &zeroArgument{t}
		}

		args[i-1] = arg
	}

	numOut := method.Type.NumOut()
	outErr := -1

	for i := numOut - 1; i >= 0; i-- {
		if method.Type.Out(i).Implements(typeOfError) {
			outErr = i
			break
		}
	}

	return &methodHandler{
		controller: reflect.ValueOf(controller),
		method:     method,
		args:       args,
		outErr:     outErr,
	}, nil
}

func (h *methodHandler) HandleRequest(c *gmvc.Context) error {
	in := make([]reflect.Value, len(h.args)+1)

	in[0] = h.controller

	for i, arg := range h.args {
		v, err := arg.Get(c)
		if err != nil {
			c.ErrorStatus(err, http.StatusBadRequest)
			return nil
		}
		in[i+1] = v
	}

	out := h.method.Func.Call(in)

	if h.outErr != -1 {
		errv := out[h.outErr]
		if !errv.IsNil() {
			return errv.Interface().(error)
		}
	}

	return nil
}
