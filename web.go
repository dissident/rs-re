package main

import (
	"fmt"
	"html"
	"net/http"
)

func serveHttp() {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.ListenAndServe(":8080", nil)
}