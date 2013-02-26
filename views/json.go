package views

import (
	"encoding/json"
	"github.com/hujh/gmvc"
)

type JsonView struct {
}

func NewJsonView() *JsonView {
	return &JsonView{}
}

func (v *JsonView) Render(c *gmvc.Context, name string, data interface{}) error {
	w := c.ResponseWriter

	h := w.Header()
	if ct := h.Get("Content-Type"); ct == "" {
		h.Set("Content-Type", "application/json")
	}

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := w.Write(d); err != nil {
		return err
	}

	return nil
}
