package gmvc

import (
	"errors"
	"net/http"
	"path"
	"strings"
)

type App struct {
	*Router
	Path         string
	Attr         Attr
	View         View
	Session      *SessionManager
	ErrorHandler ErrorHandler
}

func NewApp() *App {
	return &App{
		Path:         "/",
		Router:       NewRouter(),
		Attr:         make(Attr),
		Session:      NewSessionManager(),
		ErrorHandler: &defaultErrorHandler{},
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlpath := r.URL.Path

	if !strings.HasPrefix(urlpath, a.Path) {
		c := a.newContext(w, r)
		c.ErrorStatus(errors.New(http.StatusText(http.StatusNotFound)), http.StatusNotFound)
		return
	}

	urlpath = path.Join("/", urlpath[len(a.Path):])
	a.dispatch(w, r, urlpath)
}

func (a *App) dispatch(w http.ResponseWriter, r *http.Request, urlpath string) {
	c := a.newContext(w, r)

	ok, err := a.Router.dispatch(c, urlpath)
	if ok {
		if err != nil {
			c.Error(err)
		}
		return
	}

	status := http.StatusNotFound
	c.ErrorStatus(errors.New(http.StatusText(status)), status)
}

func (a *App) newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		app:            a,
		router:         a.Router,
		sessionManager: a.Session,
		errorHandler:   a.ErrorHandler,
		Request:        r,
		ResponseWriter: w,
		Vars:           make(PathVars),
		Attr:           make(Attr),
		View:           a.View,
		Path:           a.Path,
	}
}

type Attr map[string]interface{}
