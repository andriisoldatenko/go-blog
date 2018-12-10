package main

import (
	"net/http"
)


func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("./static"))

	mux.HandleFunc("/", index)
	mux.HandleFunc("/new/", New)
	mux.HandleFunc("/edit/", Edit)
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)
	mux.HandleFunc("/authorization-code/callback", AuthCodeCallbackHandler)

	server := &http.Server{
		Addr:     "0.0.0.0:8081",
		Handler:  mux,
	}
	server.ListenAndServe()
}