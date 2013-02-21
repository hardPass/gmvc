package gmvc

import (
	"net/http"
)

type ErrorHandler interface {
	HandleError(c *Context, err error, status int)
}

type defaultErrorHandler struct {
}

func (h *defaultErrorHandler) HandleError(c *Context, err error, status int) {
	if c.WroteHeader() {
		return
	}
	if err == nil {
		http.Error(c.ResponseWriter, "", status)
	} else {
		http.Error(c.ResponseWriter, err.Error(), status)
	}
}
