package views

import (
	"bytes"
	"fmt"
	"github.com/hujh/gmvc"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"sync"
	"text/template"
	"time"
)

const (
	wrapper = `{{$page := .}}` +
		`{{$ctx := .Context}}` +
		`{{$app := $ctx.App}}` +
		`{{$req := $ctx.Request}}` +
		`{{$res := $ctx.ResponseWriter}}` +
		`{{with .Data}}%s{{end}}`
)

type PageView struct {
	ContentType string

	mutex    sync.RWMutex
	root     string
	cache    map[string]*page
	duration time.Duration
}

func NewPageView(root string) *PageView {
	return NewPageViewWithDuration(root, -1)
}

func NewPageViewWithDuration(root string, d time.Duration) *PageView {
	return &PageView{
		ContentType: "text/html",
		root:        root,
		cache:       make(map[string]*page),
		duration:    d,
	}
}

func (v *PageView) Render(c *gmvc.Context, name string, data interface{}) error {
	w := c.ResponseWriter
	defer func() {
		c.ResponseWriter = w
	}()

	b := newBuffer(c.ResponseWriter)
	c.ResponseWriter = b

	pc := &pageContext{
		Context: c,
		Data:    data,
		view:    v,
	}

	if err := v.render(pc, name); err != nil {
		return err
	}

	h := b.Header()
	if ct := h.Get("Content-Type"); ct == "" {
		h.Set("Content-Type", v.ContentType)
	}

	if err := b.flush(); err != nil {
		return err
	}

	return nil
}

func (v *PageView) SetDuration(d time.Duration) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	v.duration = d
	for name, p := range v.cache {
		if p.expired(v.duration) {
			delete(v.cache, name)
		}
	}
}

func (v *PageView) render(pc *pageContext, name string) error {
	if reflect.ValueOf(pc.Data).IsNil() {
		pc.Data = new(empty)
	}

	t, err := v.lookup(name)
	if err != nil {
		return err
	}

	if err := t.Execute(pc.Context.ResponseWriter, pc); err != nil {
		return err
	}

	return nil
}

func (v *PageView) lookup(name string) (*template.Template, error) {
	name = filepath.Clean(name)

	v.mutex.RLock()
	p := v.cache[name]
	if p != nil && !p.expired(v.duration) {
		p.touch()
		v.mutex.RUnlock()
		return p.tpl, nil
	}
	v.mutex.RUnlock()

	t, err := v.load(name)
	if err != nil {
		return nil, err
	}
	p = &page{
		name:  name,
		tpl:   t,
		utime: time.Now(),
	}

	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.cache[name] = p
	return p.tpl, nil
}

func (v *PageView) load(name string) (*template.Template, error) {
	root, err := filepath.Abs(v.root)
	if err != nil {
		return nil, err
	}
	realpath := filepath.Join(root, name)
	b, err := ioutil.ReadFile(realpath)
	if err != nil {
		return nil, err
	}

	t := template.New(realpath)
	t = t.Funcs(pageFuncs)
	s := fmt.Sprintf(wrapper, b)
	t, err = t.Parse(s)
	if err != nil {
		return nil, err
	}

	return t, nil
}

type page struct {
	name  string
	tpl   *template.Template
	utime time.Time
}

func (p *page) expired(d time.Duration) bool {
	if d > 0 && time.Now().Sub(p.utime) > d {
		return true
	}
	return false
}

func (p *page) touch() {
	p.utime = time.Now()
}

type pageContext struct {
	Context *gmvc.Context
	Data    interface{}
	view    *PageView
}

func (c *pageContext) sub() *pageContext {
	return &pageContext{
		Context: c.Context,
		view:    c.view,
	}
}

type empty struct{}

func (e *empty) String() string {
	return ""
}

type buffer struct {
	http.ResponseWriter
	b *bytes.Buffer
}

func newBuffer(w http.ResponseWriter) *buffer {
	return &buffer{
		ResponseWriter: w,
		b:              new(bytes.Buffer),
	}
}

func (w *buffer) WriteHeader(int) {
}

func (w *buffer) Write(p []byte) (int, error) {
	return w.b.Write(p)
}

func (w *buffer) flush() error {
	_, err := w.b.WriteTo(w.ResponseWriter)
	return err
}

var pageFuncs template.FuncMap = template.FuncMap{
	"import": func(pc *pageContext, values ...interface{}) (string, error) {
		vc := len(values)
		if vc == 0 || vc > 2 {
			fmt.Errorf("wrong number of args for import: want 1 or 2 got %s", vc)
		}
		name := fmt.Sprint(values[0])
		subpc := pc.sub()
		if vc > 1 {
			subpc.Data = values[1]
		}
		return "", pc.view.render(pc, name)
	},
	"include": func(pc *pageContext, urlpath string) (string, error) {
		return "", pc.Context.Include(urlpath)
	},
}
