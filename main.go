package main

import (
	"net/http"
)


func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("./static"))

	mux.HandleFunc("/", AllPosts)
	mux.HandleFunc("/new", NewPost)
	mux.HandleFunc("/new/insert", InsertPost)
	mux.HandleFunc("/edit", EditPost)
	mux.HandleFunc("/login", LoginHandler)
	mux.HandleFunc("/logout", LogoutHandler)
	mux.HandleFunc("/authorization-code/callback", AuthCodeCallbackHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	server := &http.Server{
		Addr:     "0.0.0.0:8081",
		Handler:  mux,
	}
	server.ListenAndServe()
}