package views

import (
	"github.com/hujh/gmvc"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"
)

type TemplateView struct {
	mutex     sync.RWMutex
	root      string
	cache     map[string]*templateEntry
	cacheTime time.Duration
}

func NewTemplateView(root string) *TemplateView {
	return &TemplateView{
		root:      root,
		cache:     make(map[string]*templateEntry),
		cacheTime: 0,
	}
}

func (v *TemplateView) Render(c *gmvc.Context, name string, value interface{}) error {
	t, err := v.lookup(name)
	if err != nil {
		return err
	}

	if err := t.Execute(c.ResponseWriter, value); err != nil {
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

	s := string(b)
	t := template.New(tplpath)

	_, err = t.Parse(s)
	if err != nil {
		return nil, err
	}

	return t, nil
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
