package controllers

import (
	"github.com/hujh/gmvc"
	"io"
	"net/http"
	"reflect"
)

var (
	typeOfRequest        reflect.Type = reflect.TypeOf(new(http.Request))
	typeOfContext                     = reflect.TypeOf(new(gmvc.Context))
	typeOfPathVars                    = reflect.TypeOf(gmvc.PathVars{})
	typeOfValues                      = reflect.TypeOf(gmvc.Values{})
	typeOfMultipartForm               = reflect.TypeOf(gmvc.MultipartForm{})
	typeOfResponseWriter              = reflect.TypeOf(new(http.ResponseWriter)).Elem()
	typeOfReadCloser                  = reflect.TypeOf(new(io.ReadCloser)).Elem()
	typeOfWriter                      = reflect.TypeOf(new(io.Writer)).Elem()
)

var (
	arguments = []argument{
		&requestArgument{},
		&contextArgument{},
		&pathVarsArgument{},
		&valuesArgument{},
		&multipartFormArgument{},
		&responseWriterArgument{},
		&readCloserArgument{},
		&writerArgument{},
	}
)

type argument interface {
	Type() reflect.Type
	Get(c *gmvc.Context) (reflect.Value, error)
}

type zeroArgument struct {
	t reflect.Type
}

func (a *zeroArgument) Type() reflect.Type {
	return a.t
}

func (a *zeroArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	return reflect.Zero(a.t), nil
}

type contextArgument struct {
}

func (a *contextArgument) Type() reflect.Type {
	return typeOfContext
}

func (a *contextArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	return reflect.ValueOf(c), nil
}

type pathVarsArgument struct {
}

func (a *pathVarsArgument) Type() reflect.Type {
	return typeOfPathVars
}

func (a *pathVarsArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	return reflect.ValueOf(c.Vars), nil
}

type requestArgument struct {
}

func (a *requestArgument) Type() reflect.Type {
	return typeOfRequest
}

func (a *requestArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	return reflect.ValueOf(c.Request), nil
}

type valuesArgument struct {
}

func (a *valuesArgument) Type() reflect.Type {
	return typeOfValues
}

func (a *valuesArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	v, err := c.Form()
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	return reflect.ValueOf(v), nil
}

type multipartFormArgument struct {
}

func (a *multipartFormArgument) Type() reflect.Type {
	return typeOfMultipartForm
}

func (a *multipartFormArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	v, err := c.MultipartForm(0)
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	return reflect.ValueOf(v), nil
}

type responseWriterArgument struct {
}

func (a *responseWriterArgument) Type() reflect.Type {
	return typeOfResponseWriter
}

func (a *responseWriterArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	return reflect.ValueOf(c.ResponseWriter), nil
}

type readCloserArgument struct {
}

func (a *readCloserArgument) Type() reflect.Type {
	return typeOfReadCloser
}

func (a *readCloserArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	return reflect.ValueOf(c.Request.Body), nil
}

type writerArgument struct {
}

func (a *writerArgument) Type() reflect.Type {
	return typeOfWriter
}

func (a *writerArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	return reflect.ValueOf(c.ResponseWriter), nil
}
