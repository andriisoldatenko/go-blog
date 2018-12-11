package main

import (
	"net/http"
	"os"

	"github.com/go-http-utils/logger"
)


func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("./static"))

	mux.HandleFunc("/", AuthHandler(AllPosts))
	mux.HandleFunc("/new/", AuthHandler(NewPost))
	mux.HandleFunc("/new/insert/", AuthHandler(InsertPost))
	mux.HandleFunc("/edit/", EditPost)
	mux.HandleFunc("/login/", LoginHandler)
	mux.HandleFunc("/logout/", LogoutHandler)
	mux.HandleFunc("/authorization-code/callback/", AuthCodeCallbackHandler)
	mux.HandleFunc("/profile/", AuthHandler(ProfileHandler))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	http.ListenAndServe("0.0.0.0:8081", logger.Handler(mux, os.Stdout, logger.DevLoggerType))
}