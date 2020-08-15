package golee

import (
	"fmt"
	"net/http"
)

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP ServeHTTP
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
