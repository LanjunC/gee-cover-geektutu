package main

import (
	"fmt"
	"net/http"
)

// Engine is the uni handler for all requests
type Engine struct {}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	case "/hello":
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k , v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s", r.URL.Path)
	}
}

func main () {
	engine := new(Engine)
	http.ListenAndServe(":8080", engine)
}


