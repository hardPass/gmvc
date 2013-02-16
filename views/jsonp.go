package views

import (
	"encoding/json"
	"gmvc"
)

type JsonpView struct {
	DefaultName string
}

func NewJsonpView() *JsonpView {
	return &JsonpView{
		DefaultName: "callback",
	}
}

func (jv *JsonpView) Render(c *gmvc.Context, name string, value interface{}) error {
	w := c.ResponseWriter

	h := w.Header()
	if ct := h.Get("Content-Type"); ct == "" {
		h.Set("Content-Type", "text/javascript")
	}

	if name == "" {
		name = jv.DefaultName
	}

	if name == "" {
		name = "_"
	}

	if _, err := w.Write([]byte(name + "(")); err != nil {
		return err
	}

	d, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if _, err := w.Write(d); err != nil {
		return err
	}

	if _, err := w.Write([]byte(")")); err != nil {
		return err
	}

	return nil
}
