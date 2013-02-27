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
	http.ResponseWriter
	Request *http.Request
	Vars    PathVars
	Attrs   Attrs
	View    View
	Path    string
	Err     error

	app             *App
	parent          *Context
	request         *http.Request
	response        http.ResponseWriter
	form            Form
	multipartForm   *MultipartForm
	session         Session
	sessionProvider SessionProvider
	errorHandler    ErrorHandler
}

func (c *Context) App() *App {
	return c.app
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
	if c.multipartForm != nil {
		return c.multipartForm, nil
	}

	if c.form == nil {
		if _, err := c.Form(); err != nil {
			return nil, err
		}
	}

	if maxMemory <= 0 {
		maxMemory = defaultMaxMemory
	}
	if err := c.Request.ParseMultipartForm(maxMemory); err != nil {
		return nil, err
	}

	mf := c.Request.MultipartForm
	c.multipartForm = &MultipartForm{
		Form:  Form(mf.Value),
		Files: mf.File,
	}
	for k, v := range mf.Value {
		c.form[k] = append(c.form[k], v...)
	}

	return c.multipartForm, nil
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
		sp := c.sessionProvider
		if sp == nil {
			err = errors.New("no available session provider")
			return
		}
		s, err = sp.GetSession(c.response, c.request, create)
		if err != nil {
			return
		}
		c.session = s
	}
	return
}

func (c *Context) Include(urlpath string) error {
	urlpath = path.Join("/", urlpath)

	abspath, err := url.Parse(path.Join("/", c.Path, urlpath))
	if err != nil {
		return err
	}
	requrl := c.Request.URL.ResolveReference(abspath)

	r, err := http.NewRequest("GET", requrl.String(), nil)
	if err != nil {
		return err
	}

	for n, v := range c.Request.Header {
		r.Header[n] = v
	}

	w := newContentOnly(c.ResponseWriter)

	sc := &Context{
		Request:         r,
		ResponseWriter:  w,
		Path:            c.Path,
		Vars:            make(PathVars),
		Attrs:           make(Attrs),
		View:            c.app.View,
		app:             c.app,
		parent:          c,
		request:         c.request,
		response:        c.response,
		sessionProvider: c.sessionProvider,
		errorHandler:    c.errorHandler,
	}

	c.app.dispatch(sc, urlpath)
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
		return errors.New("no available view")
	}

	return c.View.Render(c, name, value)
}

func (c *Context) WriteString(v ...interface{}) error {
	_, err := fmt.Fprint(c.ResponseWriter, v...)
	return err
}

func (c *Context) Status(status int) {
	c.ResponseWriter.WriteHeader(status)
}

func (c *Context) Error(err error) {
	c.ErrorStatus(err, http.StatusInternalServerError)
}

func (c *Context) ErrorStatus(err error, status int) {
	c.Err = err
	h := c.errorHandler
	if h != nil {
		h.HandleError(c, err, status)
	} else {
		http.Error(c.ResponseWriter, http.StatusText(status), status)
	}
}

func (c *Context) finalize() {
	if c.parent == nil {
		if c.session != nil && c.session.Valid() {
			c.session.Save()
		}
	}
}

type Attrs map[string]interface{}

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
