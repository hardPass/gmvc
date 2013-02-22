package views

import ()

type NegotiatingView struct {
	defaultView View
	views       map[string]View
}
