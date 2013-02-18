package gmvc

type View interface {
	Render(c *Context, name string, data interface{}) error
}
