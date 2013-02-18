package gmvc

import (
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
	if !strings.HasPrefix(r.URL.Path, a.Path) {
		s := http.StatusNotFound
		if eh := a.ErrorHandler; eh != nil {
			c := a.newContext(w, r)
			c.Status(s)
		} else {
			http.Error(w, http.StatusText(s), s)
		}
		return
	}

	urlpath := path.Join("/", r.URL.Path[len(a.Path):])
	a.dispatch(w, r, urlpath)
}

func (a *App) dispatch(w http.ResponseWriter, r *http.Request, urlpath string) {
	c := a.newContext(w, r)

	if hit, err := a.Router.route(c, urlpath); hit {
		if err != nil {
			c.Error(err)
		}
		return
	}

	c.Status(http.StatusNotFound)
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
