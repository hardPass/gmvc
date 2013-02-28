package controllers

import (
	"github.com/hujh/gmvc"
	"io"
	"net/http"
	"reflect"
)

var (
	typeOfRequest       reflect.Type = reflect.TypeOf(new(http.Request))
	typeOfContext                    = reflect.TypeOf(new(gmvc.Context))
	typeOfPathVars                   = reflect.TypeOf(gmvc.PathVars{})
	typeOfValues                     = reflect.TypeOf(gmvc.Values{})
	typeOfMultipartForm              = reflect.TypeOf(gmvc.MultipartForm{})

	typeOfResponseWriter = reflect.TypeOf(new(http.ResponseWriter)).Elem()
	typeOfReadCloser     = reflect.TypeOf(new(io.ReadCloser)).Elem()
	typeOfWriter         = reflect.TypeOf(new(io.Writer)).Elem()
)

var (
	defaultArguments = []Argument{
		&RequestArgument{},
		&ContextArgument{},
		&PathVarsArgument{},
		&ValuesArgument{},
		&MultipartFormArgument{},

		&ResponseWriterArgument{},
		&ReadCloserArgument{},
		&WriterArgument{},
	}
)

type Argument interface {
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

type ContextArgument struct {
}

func (a *ContextArgument) Type() reflect.Type {
	return typeOfContext
}

func (a *ContextArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	return reflect.ValueOf(c), nil
}

type PathVarsArgument struct {
}

func (a *PathVarsArgument) Type() reflect.Type {
	return typeOfPathVars
}

func (a *PathVarsArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	return reflect.ValueOf(c.Vars), nil
}

type RequestArgument struct {
}

func (a *RequestArgument) Type() reflect.Type {
	return typeOfRequest
}

func (a *RequestArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	return reflect.ValueOf(c.Request), nil
}

type ValuesArgument struct {
}

func (a *ValuesArgument) Type() reflect.Type {
	return typeOfValues
}

func (a *ValuesArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	v, err := c.Form()
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	return reflect.ValueOf(v), nil
}

type MultipartFormArgument struct {
}

func (a *MultipartFormArgument) Type() reflect.Type {
	return typeOfMultipartForm
}

func (a *MultipartFormArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	v, err := c.MultipartForm(0)
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	return reflect.ValueOf(v), nil
}

type ResponseWriterArgument struct {
}

func (a *ResponseWriterArgument) Type() reflect.Type {
	return typeOfResponseWriter
}

func (a *ResponseWriterArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	return reflect.ValueOf(c.ResponseWriter), nil
}

type ReadCloserArgument struct {
}

func (a *ReadCloserArgument) Type() reflect.Type {
	return typeOfReadCloser
}

func (a *ReadCloserArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	return reflect.ValueOf(c.Request.Body), nil
}

type WriterArgument struct {
}

func (a *WriterArgument) Type() reflect.Type {
	return typeOfWriter
}

func (a *WriterArgument) Get(c *gmvc.Context) (reflect.Value, error) {
	return reflect.ValueOf(c.ResponseWriter), nil
}
