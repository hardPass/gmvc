package views

import (
	"encoding/json"
	"github.com/hujh/gmvc"
)

type JsonpView struct {
	DefaultName string
}

func NewJsonpView() *JsonpView {
	return &JsonpView{
		DefaultName: "callback",
	}
}

func (v *JsonpView) Render(c *gmvc.Context, name string, data interface{}) error {
	w := c.ResponseWriter

	if !c.WroteHeader() {
		h := w.Header()
		if ct := h.Get("Content-Type"); ct == "" {
			h.Set("Content-Type", "text/javascript")
		}
	}

	fn := name
	if fn == "" {
		fn = v.DefaultName
		if fn == "" {
			fn = "_"
		}
	}

	if _, err := w.Write([]byte(fn + "(")); err != nil {
		return err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	if _, err := w.Write([]byte(")")); err != nil {
		return err
	}

	return nil
}
