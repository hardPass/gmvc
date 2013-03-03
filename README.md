# gmvc

web framework for golang

## Installation

To install gmvc:

```
go get github.com/hujh/gmvc
```

## Quick Start
    
```go
package main

import (
	"log"
	"net/http"
	"github.com/hujh/gmvc"
)

func main() {
	app := gmvc.NewApp()

	app.HandleFunc("/", func(c *gmvc.Context) error {
		return c.WriteString("hello world!")
	})

	log.Fatal(s.ListenAndServe(":8080", app))
}
```

## Routing

A gmvc Router matches incoming requests and calls filters and handler for the pattern that matches the URL.

Pattern syntax: [HttpMethods] <PathPattern>

HttpMethods: a list of HTTP Method, separated by ',' (optional)
PathPattern: a pattern of uri path. use {pathvar:regexp} to extract vars from URIs (require).

Here's a sample:
```
GET      /index          // simple
/index                   // simple (omit http methods)
GET,POST /user/profile   // multi http method
*        /hello          // all http method
GET      /user/{id:\d+}  // extract a numberic pathvar 'id'
GET      /book/**/index  // wildcard
```

Handler
Handler should accept and service request, use *Context to read and reply data.
 
```go
type Handler interface {
	HandleRequest(*Context) error
}

...
router.HandleFunc("GET /users", func(c *gmvc.Context) error {
	// ...
})
```

Filter

Filter perform logic either before match handler or after an handler serviced.

```go
type Filter interface {
	DoFilter(fc *FilterContext) error
}

...
router.FilterFunc("/admin/**", func(fc *FilterContext){
	//...
	return fc.Next()
})
```

## Context
Once receiving a request, gmvc will wrap the http.ResponseWriter and http.Request as a context, It provides useful methods to store or output data to client.

```go
type Context struct {
	http.ResponseWriter
	Request *http.Request
	Vars    PathVars
	Attrs   Attrs
	View    View
	Path    string
	Err     error
}
```

## Input Values

PathVars:
PathVars extracts vars from pattern which define var and assgin to the FilterContext or Context

```go
app.HandleFunc("/{a}/{b}", func(c *gmvc.Context) error {
	return c.WriteString(c.Vars[a], c.Vars[b])
})

// "/book/123" --> "book 123"
```

Form Values

Values: Values contains the parsed form data, including both the URL field's query parameters and the POST or PUT form data (like url.Values), it supports common convert from string to sespecial type.

MultipartForm: MultipartForm is the parsed multipart form, including file uploads.


example
```go
app.HandleFunc("/book", func(c *gmvc.Context) error {
	f, err := c.Form()
	if err != nil {
		return err
	}
	return c.WriteString(f.int("year"))
})

// "/book?year=1990" --> "1990"
```

### Session
gmvc provides SessionProvider and Session interface to support session. Sometimes user must implements them to satisfy the requirements.

```go
type SessionProvider interface {
	GetSession(w http.ResponseWriter, r *http.Request, create bool) (Session, error)
}

type Session interface {
	Id() string

	Valid() bool
	Invalidate() error

	Save() error

	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Del(key string) error
}
```

Get the session from Context:
```go
session, err := c.Session(true)
```

There is a memory session implements.

```go
import "github.com/hujh/gmvc/sessions"

...
app.SessionProvider = sessions.NewMemoryProvider(30 * time.Minute)

```

## View

View is a component that render the result, it is a interface type:
```go
type View interface {
	Render(c *Context, name string, data interface{}) error
}
```

To use it before, you must assign to the app or context like:
```go
app.View = view     // app scope
```
or
```go
context.View = view // request scope
```

and then render
```go
context.Render(name, data)
```

install some useful buildin Views
```go
import "github.com/hujh/gmvc/views"
```

## Controllers

controllers is a helper which use to maps pattern to method for the Router

import before
```go
import "github.com/hujh/gmvc/controllers"
```

user controller must implements Controller interface:
```go
type Controller interface {
	RequestMapping() string
}
```
RequestMapping returns the string mapping which separated by line end. Every not empty line maps a pattern to a controller method. syntax: 
```
[HttpMethods] <PathPattern> <ControllerMethod>
...
```

example:
```go
controllers.Register(app.Router, "/user", &UserController{})

...

type UserController struct {
}

func (uc *UserController) RequestMapping() string {
	return `
	GET   /{id}   GET
	POST  /{id}   Create
	GET   /list   List
	`
}

func (uc *UserController) Get(c *gmvc.Context) error { ... }
func (uc *UserController) Create(c *gmvc.Context) error { ... }
func (uc *UserController) List(c *gmvc.Context, w http.ResponseWriter) error { ... }
```

if user controller method defines more arguments (like Request, ResponseWriter), the arguments will be auto inject.

Support arguments:
```
*http.Request
*gmvc.Context
*gmvc.PathVars
*gmvc.Values
*gmvc.MultipartForm
http.ResponseWriter
io.ReadCloser
io.Writer
```

TO BE CONTINUED...

