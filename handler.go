package gmvc

type Handler interface {
	HandleHTTP(c *Context) error
}

type HandlerFunc func(*Context) error

func (f HandlerFunc) HandleHTTP(c *Context) error {
	return f(c)
}
