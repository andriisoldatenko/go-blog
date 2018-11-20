---
layout: blog_post
title: "Go web app in 20 minutes"
author: a_soldatenko
description: "This tutorial explains how to build a web application using Go in less than 20 minutes."
tags: [go, angular, crud]
---

## [DRAFT] Go web app in 20 minutes

## Go and web applications
Nowadays web is the most used kind of software applications. Each day we sniff internet using laptops 
and mobile devices, we play in games and watch TV shows.
HTTP is the basis for communicating in web.


## Prerequisites
Before we dig into building web app in Go, please make sure you have installed latest Go in your 
operating system. To check current installed version you can run:
```bash
$ go version
go version go1.11 darwin/amd64
```

## Few words about go modules.

## Project layout
```bash
tree -L 3
go-blog git:(master) ✗ tree -L 3
.
├── README.md
├── TODO.md
├── examples
│   └── hello_web
│       └── main.go
├── go-blog
├── go.mod
├── go.sum
├── main.go
├── main_test.go
├── serve_static.go
└── static
    └── hello.txt
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


## Let's build simple and modern blog platform

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

## Nginx with go server?
Do I really need nginx behind go server?
It's a good question :) 

## Few words about debugging

## Code reloading on change
