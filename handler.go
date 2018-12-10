package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var tmpl = template.Must(template.ParseGlob("templates/*"))


// Get all Authors and render template with Authors
func Index(w http.ResponseWriter, r *http.Request) {
	db := DBConn()
	var authors []Author
	err := db.Model(&authors).Select()
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(w, "Index", r)
	defer db.Close()
}

// Return new blog Post html form on GET
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

// Create new author post using form submit
func Insert(w http.ResponseWriter, r *http.Request) {
	db := DBConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		user1 := &Author{
			Name: name,
		}
		err := db.Insert(user1)
		if err != nil {
			panic(err)
		}
		log.Println("Create Post: Name: " + name)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

// Update author details
func Edit(w http.ResponseWriter, r *http.Request) {
	db := DBConn()
	nId, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	// Select user by primary key.
	user := &Author{Id: nId}
	err = db.Select(user)
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(w, "Edit", user)
	defer db.Close()
}