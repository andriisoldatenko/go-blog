---
layout: blog_post
title: "Go web app in 20 minutes"
author: a_soldatenko
description: "This tutorial explains how to build a web application using Go in less than 20 minutes."
tags: [go, html, blog]
---

# Go web app in 20 minutes

## Go and web applications
Nowadays web is the most used kind of software applications. Each day we sniff internet using laptops 
and mobile devices, we play in games and watch TV shows. HTTP is the basis for communicating in web.
The first version of HTTP has been created by Tim Berners-Lee in 1990, was a simple protocol created 
to help adoption of the World Wide Web. Everything that you see on a web page is transported through 
this simple text-based protocol. HTTP is a stateless, text-based, request-response protocol that uses 
the client-server computing model. Go as many other programming languages has all needed out-of-box 
for supporting http through [`net/http`](https://golang.org/pkg/net/http/) package.

## Prerequisites
Before we dig into building web app in Go, please make sure you have installed latest Go in your 
operating system. To check current installed version you can run:
```bash
$ go version
go version go1.12 darwin/amd64
```
For more details how to install fo indifferent operation systems please visit 
[go install](https://golang.org/doc/install) page.

## Few words about go modules.
Since Go1.11 has been released we finally can see go official way for dependency management.
Hovever, modules are an experimental feature in Go 1.11, go community are going to finalize it in
Go 1.12. More details in [Go 1.11 Modules](https://github.com/golang/go/wiki/Modules).


## Project layout
```bash
tree -L 3
go-blog git:(master) ✗ tree -L 3
.
4 directories, 10 files
├── Makefile
├── README.md
├── blog
├── db.go
├── examples
├── go.mod
├── go.sum
├── handler.go
├── main.go
├── session.go
├── sql
├── static
├── templates
└── utils.go
```

## Writing your first Go web app
I prefer to start everything by example:
```go
package main

import (
	"fmt"
	"net/http"
)

func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello %s!", request.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
```
A _handler_ receives and processes the HTTP request sent from the client. It also calls the 
template engine to generate the HTML and finally bundles data into the HTTP response to be sent 
back to the client.


## Hello World! 
Let's build simple and modern blog platform, but before we start, let's review basic `server.go` example:
Open you favorite text editor, create file and paste listing belows: 
```go
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
```
Now we are good to run it:
```bash
$ go run server.go
```
and now try to open you favorite browser and navigate to `http://localhost:8081`.

## Templates


## Database
For storage in selected cool DB PostgreSQL, let's pull latest version using Docker:
```bash
$ docker run --name postgresql -itd --restart always \
    --publish 5432:5432 \
    --volume /Users/andrii/docker/postgresql:/var/lib/postgresql \
    --env 'DB_USER=blog' --env 'DB_PASS=blog_secret_password' \
    --env 'DB_NAME=blog_db' \
    sameersbn/postgresql:10
```

Now we can check, that database is ready to accept connections:
```bash
docker logs -f postgresql
# few missing lines
2018-11-21 16:20:04.714 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
2018-11-21 16:20:04.714 UTC [1] LOG:  listening on IPv6 address "::", port 5432
2018-11-21 16:20:04.718 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
2018-11-21 16:20:04.801 UTC [1105] LOG:  database system was shut down at 2018-11-21 16:20:04 UTC
2018-11-21 16:20:04.834 UTC [1] LOG:  database system is ready to accept connections
```
And also it always make sense to try to connect using `psql` command line tool:
```bash
 psql -h localhost -U blog -d blog_db -p 5432
Password for user blog:
psql (10.5, server 10.4 (Ubuntu 10.4-2.pgdg18.04+1))
Type "help" for help.

blog_db=> \dt;
Did not find any relations.
blog_db=>
```
As you can see, there is no any tables in our just created database `blog`. In real world application, you need
to create  
## Models
In Go typically we are using type `struct` for mapping database tables. In our blog app, we need to models or
2 structs `User`:

```
type User struct {
	Id          int
	Name        string
	Email       string
	Password    string
	CreatedAt   time.Time
	PublishedAt time.Time
}
```
 and `Post`, and use has many posts relationship.

```
type Post struct {
	Id          int
	Author      string
	Title       string
	Text        string
	UserId      User
	CreatedAt   time.Time
	PublishedAt time.Time
}
```

# Authentication using Okta
To divide our readers and writers of blog posts, it make sense that writers must have account and
authorized before editing posts, but readers can read without any credentials.
Let's add Okta login for blog authors:
```go

```

First you need to create Okta Platform account and verify you email.
After successfully created account, you need to copy `.env.example` to `.env` and replace with your values:

```bash
CLIENT_ID=
CLIENT_SECRET=
ISSUER=https://{yourOktaDomain}.com/oauth2/default
```




## Fin
🎉 🎉 🎉 Congratulations! You build your first blog in go. 


## Some useful references
-  



A quick introduction to HTTP
Request-response cycle