package gmvc

type Handler interface {
	HandleRequest(*Context) error
}

type HandlerFunc func(*Context) error

func (f HandlerFunc) HandleRequest(c *Context) error {
	return f(c)
}
