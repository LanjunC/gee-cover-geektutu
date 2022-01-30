package main

import (
	"fmt"
	"log"
	"net/http"
)

/**
Go语言内置的 net/http库，封装了HTTP网络编程的基础的接口
 */
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handler echoes r.URL.Path
func indexHandler( w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k , v)
	}
}

