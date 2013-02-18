package views

import (
	"bytes"
	"fmt"
	"github.com/hujh/gmvc"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	"time"
)

const (
	vars   string = `{{$ctx := .C}}{{$app := $ctx.App}}{{$req := $ctx.Request}}{{$res := $ctx.ResponseWriter}}`
	layout        = `%s{{with .D}}%s{{end}}`
)

var funcMap template.FuncMap = template.FuncMap{
	"include": func(c *gmvc.Context, urlpath string) (string, error) {
		err := c.Include(urlpath)
		return "", err
	},
}

type TemplateView struct {
	ContentType string

	mutex     sync.RWMutex
	root      string
	cache     map[string]*templateEntry
	cacheTime time.Duration
}

func NewTemplateView(root string) *TemplateView {
	return &TemplateView{
		ContentType: "text/html",
		root:        root,
		cache:       make(map[string]*templateEntry),
		cacheTime:   0,
	}
}

func (v *TemplateView) Render(c *gmvc.Context, name string, data interface{}) error {
	t, err := v.lookup(name)
	if err != nil {
		return err
	}

	w := c.ResponseWriter
	defer func() {
		c.ResponseWriter = w
	}()

	bw := newBufferWriter(c.ResponseWriter)
	c.ResponseWriter = bw

	tdata := &templateData{
		C: c,
		D: data,
	}

	if err := t.Execute(bw, tdata); err != nil {
		return err
	}

	h := c.ResponseWriter.Header()
	if ct := h.Get("Content-Type"); ct == "" {
		h.Set("Content-Type", v.ContentType)
	}

	if err := bw.flush(); err != nil {
		return err
	}

	return nil
}

func (v *TemplateView) SetCacheTime(ct time.Duration) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	v.cacheTime = ct
	for name, entry := range v.cache {
		if entry.expired(v.cacheTime) {
			delete(v.cache, name)
		}
	}
}

func (v *TemplateView) lookup(name string) (*template.Template, error) {
	v.mutex.RLock()

	name = filepath.Clean(name)
	entry := v.cache[name]

	if entry != nil && !entry.expired(v.cacheTime) {
		entry.touch()
		v.mutex.RUnlock()
		return entry.tpl, nil
	}

	v.mutex.RUnlock()

	v.mutex.Lock()
	defer v.mutex.Unlock()

	t, err := v.load(name)
	if err != nil {
		return nil, err
	}

	entry = newTemplateEntry(name, t)
	v.cache[name] = entry

	return entry.tpl, nil
}

func (v *TemplateView) load(name string) (*template.Template, error) {
	root, err := filepath.Abs(v.root)
	if err != nil {
		return nil, err
	}

	tplpath := filepath.Join(root, name)

	b, err := ioutil.ReadFile(tplpath)
	if err != nil {
		return nil, err
	}

	s := fmt.Sprintf(layout, vars, b)
	t := template.New(tplpath)
	t = t.Funcs(funcMap)

	_, err = t.Parse(s)
	if err != nil {
		return nil, err
	}

	return t, nil
}

type bufferWriter struct {
	http.ResponseWriter
	b *bytes.Buffer
}

func newBufferWriter(w http.ResponseWriter) *bufferWriter {
	return &bufferWriter{
		ResponseWriter: w,
		b:              new(bytes.Buffer),
	}
}

func (w *bufferWriter) WriteHeader(int) {
}

func (w *bufferWriter) Write(p []byte) (int, error) {
	return w.b.Write(p)
}

func (w *bufferWriter) flush() error {
	_, err := w.b.WriteTo(w.ResponseWriter)
	return err
}

type templateData struct {
	C *gmvc.Context
	D interface{}
}

type templateEntry struct {
	name  string
	tpl   *template.Template
	utime time.Time
}

func newTemplateEntry(name string, tpl *template.Template) *templateEntry {
	return &templateEntry{
		name:  name,
		tpl:   tpl,
		utime: time.Now(),
	}
}

func (e *templateEntry) expired(t time.Duration) bool {
	if t == 0 {
		return false
	}
	if time.Now().Sub(e.utime) > t {
		return true
	}
	return false
}

func (e *templateEntry) touch() {
	e.utime = time.Now()
}
