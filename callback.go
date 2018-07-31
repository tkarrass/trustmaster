package trustmaster

import (
	"net/http"
	"fmt"
)

type AuthEvent struct {
	Code string
	Tag string
}

type CallbackHandler struct {
	Events chan AuthEvent
	Redirect string  // redirect to this target after reading the code and tag
	Error string  // redirect here, if any error occured. will redirect to Redirect if empty
}

// Creates a http callback handler which reads the code and additional tag
// from the request and forwards to a given redirect target.
func (handler CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get code and tag
	code := r.URL.Query().Get("code")
	tag := r.URL.Query().Get("state")
	target := handler.Redirect
	if code == "" || tag == "" {
		if handler.Error != "" {
			target = handler.Error
		}
	} else {
		handler.Events <- AuthEvent{code, tag}
	}

	// redirect
	w.Header().Set("Location", target)
	w.WriteHeader(307)
	fmt.Fprintf(w, "<html><body><h3>moved</h3></body></html>")
}