package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var Db *sql.DB

type User struct {
	Id          int64
	Name        string
	Email       string
	Password    string
	CreatedAt   time.Time
	PublishedAt time.Time
}


type Post struct {
	Id          int64
	Author      *User
	Title       string
	Text        string
	CreatedAt   time.Time
	PublishedAt time.Time
}

func init() {
	var err error
	Db, err = sql.Open("postgres", "dbname=blog sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return
}