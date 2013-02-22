package views

import (
	"github.com/hujh/gmvc"
)

// TODO: Content Negotiating support
type NegotiatingView struct {
	defaultView gmvc.View
	views       map[string]gmvc.View
}
