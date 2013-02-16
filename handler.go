package gmvc

import ()

type Handler interface {
	Serve(c *Context) error
}

type HandlerFunc func(*Context) error

func (f HandlerFunc) Serve(c *Context) error {
	return f(c)
}
