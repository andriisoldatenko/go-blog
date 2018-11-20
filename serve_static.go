package main

import (
	"fmt"
	"net/http"
)


func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello %s!", request.URL.Path[1:])
}

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", handler)
	server := &http.Server{
		Addr:     "0.0.0.0:8081",
		Handler:  mux,
	}
	server.ListenAndServe()
}