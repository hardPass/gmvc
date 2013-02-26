package gmvc

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

type Context struct {
	Request *http.Request
	http.ResponseWriter
	Vars PathVars
	Attr Attr
	View View
	Path string
	Err  error

	app             *App
	parent          *Context
	request         *http.Request
	response        *response
	form            Form
	multipartForm   *MultipartForm
	session         Session
	sessionProvider SessionProvider
	errorHandler    ErrorHandler
}

func (c *Context) App() *App {
	return c.app
}

func (c *Context) WroteHeader() bool {
	return c.response.wroteHeader
}

func (c *Context) Form() (Form, error) {
	if c.form != nil {
		return c.form, nil
	}

	if err := c.Request.ParseForm(); err != nil {
		return nil, err
	}

	c.form = Form(c.Request.Form)
	return c.form, nil
}

func (c *Context) MultipartForm(maxMemory int64) (*MultipartForm, error) {
	f := c.multipartForm
	if f != nil {
		return f, nil
	}

	if maxMemory <= 0 {
		maxMemory = defaultMaxMemory
	}

	if err := c.Request.ParseMultipartForm(maxMemory); err != nil {
		return nil, err
	}
	return (*MultipartForm)(c.Request.MultipartForm), nil
}

func (c *Context) Session(create bool) (s Session, err error) {
	if c.parent != nil {
		return c.parent.Session(create)
	}

	s = c.session
	if s != nil && !s.Valid() {
		s = nil
	}

	if s == nil {
		s, err = c.sessionProvider.GetSession(c.response, c.request, create)
		if err != nil {
			return
		}
		c.session = s
	}

	return
}

func (c *Context) Include(urlpath string) error {
	r := c.Request
	w := c.ResponseWriter

	iw := newContentOnly(w)
	urlpath = path.Join("/", urlpath)

	iru, err := url.Parse(path.Join("/", c.Path, urlpath))
	if err != nil {
		return err
	}

	iu := r.URL.ResolveReference(iru)
	ir, err := http.NewRequest("GET", iu.String(), nil)
	if err != nil {
		return err
	}

	for n, v := range r.Header {
		ir.Header[n] = v
	}

	ic := &Context{
		Request:         ir,
		ResponseWriter:  iw,
		Path:            c.Path,
		Vars:            make(PathVars),
		Attr:            make(Attr),
		View:            c.app.View,
		app:             c.app,
		parent:          c,
		request:         c.request,
		response:        c.response,
		sessionProvider: c.sessionProvider,
		errorHandler:    c.errorHandler,
	}

	c.app.dispatch(ic, urlpath)
	return nil
}

func (c *Context) Redirect(urlstr string, code int) error {
	u, err := url.Parse(urlstr)
	if err != nil {
		return err
	}

	if !u.IsAbs() {
		urlstr = path.Join(path.Join("/", c.Path, urlstr))
	}

	http.Redirect(c.ResponseWriter, c.Request, urlstr, code)
	return nil
}

func (c *Context) Render(name string, value interface{}) error {
	if c.View == nil {
		return fmt.Errorf("not found view: %s", name)
	}

	return c.View.Render(c, name, value)
}

func (c *Context) WriteString(v ...interface{}) error {
	_, err := c.ResponseWriter.Write([]byte(fmt.Sprint(v...)))
	return err
}

func (c *Context) Status(status int) {
	if status >= 400 {
		c.ErrorStatus(errors.New(http.StatusText(status)), status)
	} else {
		c.ResponseWriter.WriteHeader(status)
	}
}

func (c *Context) Error(err error) {
	c.ErrorStatus(err, http.StatusInternalServerError)
}

func (c *Context) ErrorStatus(err error, status int) {
	c.Err = err
	h := c.errorHandler
	if h != nil {
		h.HandleError(c, err, status)
	}
}

func (c *Context) finalize() {
	if c.parent == nil {
		if c.session != nil && c.session.Valid() {
			c.session.Save()
		}
	}
}

type contentOnly struct {
	h http.Header
	w http.ResponseWriter
}

func newContentOnly(w http.ResponseWriter) *contentOnly {
	return &contentOnly{
		h: make(http.Header),
		w: w,
	}
}

func (c *contentOnly) Header() http.Header {
	return c.h
}

func (c *contentOnly) WriteHeader(int) {
}

func (c *contentOnly) Write(p []byte) (int, error) {
	return c.w.Write(p)
}
