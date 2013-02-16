package gmvc

import (
	"net/http"
)

type ErrorHandler interface {
	HandleError(c *Context, status int, err error)
}

type defaultErrorHandler struct {
}

func (h *defaultErrorHandler) HandleError(c *Context, status int, err error) {
	w := c.ResponseWriter

	if err == nil {
		http.Error(w, "", status)
	} else {
		http.Error(w, err.Error(), status)
	}
}
