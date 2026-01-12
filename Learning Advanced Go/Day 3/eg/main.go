package main

import (
	"fmt"
	"net/http"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("user")
	if name == "" {
		name = "Guest"
	}
	fmt.Fprintf(w, "Welcome, %s!", name)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/greet", welcome)
	http.ListenAndServe(":8080", mux)
}
