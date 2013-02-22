package gmvc

import (
	"net/http"
)

// TODO: middleware support
type Middleware interface {
	ProcessRequest(w http.ResponseWriter, r *http.Request)
}
