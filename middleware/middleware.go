package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Add takes an httpHandler and calls all specified handlers for adding middleware
func Add(h httprouter.Handle, middleware ...func(httprouter.Handle) httprouter.Handle) httprouter.Handle {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

// SetHeaders adds required headers to all requests
func SetHeaders(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		h(w, r, ps)
	}
}
