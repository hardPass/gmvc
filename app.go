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
	res := newResponse(w)

	return &Context{
		Request:        r,
		ResponseWriter: res,
		Path:           a.Path,
		Vars:           make(PathVars),
		Attr:           make(Attr),
		View:           a.View,
		app:            a,
		response:       res,
		sessionManager: a.Session,
		errorHandler:   a.ErrorHandler,
	}
}

type Attr map[string]interface{}

type response struct {
	http.ResponseWriter
	wroteHeader bool
}

func newResponse(w http.ResponseWriter) *response {
	return &response{ResponseWriter: w}
}

func (r *response) WriteHeader(status int) {
	if !r.wroteHeader {
		r.wroteHeader = true
	}
	r.ResponseWriter.WriteHeader(status)
}

func (r *response) Write(p []byte) (int, error) {
	if !r.wroteHeader {
		r.wroteHeader = true
	}
	return r.ResponseWriter.Write(p)
}
