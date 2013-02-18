package views

import (
	"encoding/xml"
	"github.com/hujh/gmvc"
)

type XmlView struct {
}

func NewXmlView() *XmlView {
	return &XmlView{}
}

func (v *XmlView) Render(c *gmvc.Context, name string, data interface{}) error {

	w := c.ResponseWriter
	h := w.Header()

	if ct := h.Get("Content-Type"); ct == "" {
		h.Set("Content-Type", "text/xml")
	}

	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
	d, err := xml.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := w.Write(d); err != nil {
		return err
	}

	return nil
}
