package controllers

import (
	"gmvc"
	"net/http"
	"reflect"
)

var (
	typeOfError = reflect.TypeOf(new(error)).Elem()
)

type methodHandler struct {
	controller reflect.Value
	method     reflect.Method
	args       []Argument
	outErr     int
}

func newMethodHandler(controller interface{}, method reflect.Method, arguments []Argument) (*methodHandler, error) {

	numIn := method.Type.NumIn()
	args := make([]Argument, numIn-1)

	for i := 1; i < numIn; i++ {
		var arg Argument
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

	h := &methodHandler{
		controller: reflect.ValueOf(controller),
		method:     method,
		args:       args,
		outErr:     outErr,
	}

	return h, nil
}

func (h *methodHandler) Serve(c *gmvc.Context) error {
	in := make([]reflect.Value, len(h.args)+1)

	in[0] = h.controller

	for i, arg := range h.args {
		v, err := arg.Get(c)
		if err != nil {
			c.ErrorStatus(err, http.StatusBadRequest)
		}
		in[i+1] = v
	}

	out := h.method.Func.Call(in)

	if h.outErr != -1 {
		errv := out[h.outErr]
		if !errv.IsNil() {
			err := errv.Interface().(error)
			return err
		}
	}

	return nil
}
