package main

import (
	"github.com/gorilla/context"
	"html/template"
	"log"
	"net/http"
	"strconv"
)


// Get all blog posts and render template
func AllPosts(w http.ResponseWriter, r *http.Request) {
	db := DBConn()
	var posts []Post
	err := db.Model(&posts).Select()
	if err != nil {
		panic(err)
	}
	t, _ := template.ParseFiles("templates/layout.html", "templates/index.html")
	data := context.Get(r, "data")
	extra := struct {
		CustomData
		Posts []Post
	}{CustomData: data.(CustomData), Posts: posts}
	t.Execute(w, extra)
	defer db.Close()
}

// Return new blog Post html form on GET
func NewPost(w http.ResponseWriter, r *http.Request) {
	data := context.Get(r, "data")
	t, _ := template.ParseFiles("templates/layout.html", "templates/new.html")
	t.Execute(w, data)
}

// Create new blog post post using form submit
func InsertPost(w http.ResponseWriter, r *http.Request) {
	db := DBConn()
	if r.Method == "POST" {
		title := r.FormValue("title")
		content := r.FormValue("content")
		email := context.Get(r, "email").(string)
		post1 := &Post{
			Title: title,
			Content: content,
			AuthorEmail: email,
		}
		err := db.Insert(post1)
		if err != nil {
			panic(err)
		}
		log.Println("Create Blog Post: Title: " + title)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

// Update post details
func EditPost(w http.ResponseWriter, r *http.Request) {
	db := DBConn()
	nId, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	// Select user by primary key.
	post := &Post{Id: nId}
	err = db.Select(post)
	if err != nil {
		panic(err)
	}
	t, _ := template.ParseFiles("templates/layout.html", "templates/edit.html")
	t.Execute(w, post)
	defer db.Close()
}