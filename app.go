package gmvc

import (
	"errors"
	"net/http"
	"path"
	"strings"
	"sync"
)

type App struct {
	*Router
	Path            string
	Attrs           *AppAttrs
	View            View
	SessionProvider SessionProvider
	ErrorHandler    ErrorHandler
}

func NewApp() *App {
	return &App{
		Path:         "/",
		Router:       NewRouter(),
		Attrs:        &AppAttrs{},
		ErrorHandler: &defaultErrorHandler{},
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := a.buildContext(w, r)

	if !strings.HasPrefix(r.URL.Path, a.Path) {
		defer c.finalize()
		errorStatus(c, http.StatusNotFound)
		return
	}

	urlpath := path.Join("/", r.URL.Path[len(a.Path):])
	a.dispatch(c, urlpath)
}

func (a *App) dispatch(c *Context, urlpath string) {
	defer c.finalize()

	if hit, err := a.Router.route(c, urlpath); hit {
		if err != nil {
			c.Error(err)
		}
		return
	}

	errorStatus(c, http.StatusNotFound)
}

func (a *App) buildContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Request:         r,
		ResponseWriter:  w,
		Path:            a.Path,
		Vars:            make(PathVars),
		Attrs:           make(Attrs),
		View:            a.View,
		app:             a,
		request:         r,
		response:        w,
		sessionProvider: a.SessionProvider,
		errorHandler:    a.ErrorHandler,
	}
}

type AppAttrs struct {
	mutex  sync.RWMutex
	values map[string]interface{}
}

func (a *AppAttrs) Set(key string, value interface{}) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.values[key] = value
}

func (a *AppAttrs) Get(key string) interface{} {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.values[key]
}

func (a *AppAttrs) Del(key string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	delete(a.values, key)
}

func errorStatus(c *Context, status int) {
	err := errors.New(http.StatusText(status))
	c.ErrorStatus(err, status)
}
