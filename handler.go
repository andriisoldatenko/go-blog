package main

import (
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
	t.Execute(w, posts)
	defer db.Close()
}

// Return new blog Post html form on GET
func NewPost(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/layout.html", "templates/new.html")
	t.Execute(w, nil)
}

// Create new blog post post using form submit
func InsertPost(w http.ResponseWriter, r *http.Request) {
	db := DBConn()
	if r.Method == "POST" {
		title := r.FormValue("title")
		post1 := &Post{
			Title: title,
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