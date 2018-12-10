---
layout: blog_post
title: "Go web app in 20 minutes"
author: a_soldatenko
description: "This tutorial explains how to build a web application using Go in less than 20 minutes."
tags: [go, html, blog]
---

# Go web app in 20 minutes

## Go and web applications

Nowadays, web applications are the most used kind of software applications. Each day we surf the internet using laptops and mobile devices, we play in games and watch TV shows. HTTP is the basis for communicating in web. The first version of HTTP (created by Tim Berners-Lee in 1990) was a simple protocol created to help the adoption of the World Wide Web. Everything you see on a web page is transported through this simple text-based protocol. HTTP is a stateless, text-based, request-response protocol that uses the client-server computing model. Go as many other programming languages has all needed out-of-box 
for supporting http through [`net/http`](https://golang.org/pkg/net/http/) package.

## Prerequisites

Before we dig into building web app in Go, please make sure you have installed latest Go in your 
operating system. To check current installed version you can run:

```bash
$ go version
go version go1.11.2 darwin/amd64
```

## Writing your first Go web app

I prefer to start everything by example, let's create first file `server.go` using you favorite text editor:

```bash
$ cat server.go
```

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

Now we can build and run our first go web app!:

```bash
go run server.go
```

Note: `go run` - build and run your go program in the same time.

Now if you open you browser and go to `http://localhost:8081/Okta`, you can see:

```bash
$ curl http://localhost:8081/okta
Hello okta!
```

## Go modules.

Since Go 1.11 has been released we finally can see go official way for dependency management.
However, modules are an experimental feature in Go 1.11, go community are going to finalize it in
Go 1.12. More details in [Go 1.11 Modules](https://github.com/golang/go/wiki/Modules).

### Create project folder

Quick start with Go Modules:

```bash
$ mkdir -p $HOME/work/go-blog
$ cd $HOME/work/go-blog
$ go mod init github.com/you/go-blog
$ ls -la
-rw-r--r--    1 andrii  staff    30 Dec 10 14:02 go.mod
-rw-r--r--    1 andrii  staff  1374 Dec 10 14:01 main.go
$ cat go.mod
module github.com/you/go-blog
```


Now we can write some code with dependencies.

## Let's build a simple and modern blog platform

## Models

In Go typically we are using type `struct` for mapping database tables. In our blog app, we need two models or structs `Author`:

```
type Author struct {
	Id          int64
	Name        string
	Email       string
	Password    string
}
```

and `Post`, and use has many posts relationship.

```
type Post struct {
	Id          int64
	Author      *Author
	AuthorId    int64
	Title       string
	Content     string
}
```

And now we can create schema and connect to DB:
Let's create `db.go` file where all database related code are placed in `go-blog` folder:
```bash
$ cat db.go
```

```go
package main

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type Author struct {
	Id          int64
	Name        string
	Email       string
	Password    string
}

func (u Author) String() string {
	return fmt.Sprintf("Author<%d %s %v>", u.Id, u.Name, u.Email)
}

type Post struct {
	Id       int64
	Title    string
	AuthorId int64
	Author   *Author
}

func (s Post) String() string {
	return fmt.Sprintf("Post<%d %s %s>", s.Id, s.Title, s.Author)
}


func DBConn() (db *pg.DB) {
	db = pg.Connect(&pg.Options{
		Database: "blog_db",
		User: "blog",
		Password: "blog_secret_password",
	})
	return db
}

func init() {
	db := DBConn()
	createSchema(db)
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Author)(nil), (*Post)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
``` 

Before we try to build model related code, we should return back to Go modules, and `go-pg` as dependency:
It's very easy, as soon you run `go build` or `go run`, go will do it automatically:

```bash
go run main.go
go: finding github.com/go-pg/pg/orm latest
go: downloading github.com/go-pg/pg v6.15.1+incompatible
go: finding github.com/jinzhu/inflection latest
...
```

Now you can see another issue with we forgot to add / install our database and by default we are trying to reach localhost:5432 PostgreSQL instance:
```
panic: dial tcp [::1]:5432: connect: connection refused
```
It means we need to configure our database.

## Database

For storage in selected cool DB PostgreSQL as the database, let's pull latest version using Docker:

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
And also it always makes sense to try to connect using `psql` command line tool:

```bash
 psql -h localhost -U blog -d blog_db -p 5432
Password for user blog:
psql (10.5, server 10.4 (Ubuntu 10.4-2.pgdg18.04+1))
Type "help" for help.

blog_db=> \dt;
Did not find any relations.
blog_db=>
```

As you can see, there is no any tables in our just created database `blog`. In real world application, you need to create  

## Integrate DB with models

Now we can edit our `main.go` to replace with DB connection settings:

```go
db := pg.Connect(&pg.Options{
    Database: "blog_db",
    User: "blog",
    Password: "blog_secret_password",
})
```

If you run `go run main.go` you can see that we generated tables and insert `Author` and `Post`:
```bash
Author<1 admin >
[Author<1 admin >]
Post<1 Cool story Author<1 admin >>
```

## Handlers or views
A view function, or view for short, is simply a Go function that takes a Web request and write to response object. In our example this response can be the HTML contents of a web page, or a redirect, or a 404 error, or an XML document, or an image . . . or anything, really. The view itself contains whatever arbitrary logic is necessary to return that response. This code can live anywhere you want, as long as it’s on your Python path. There’s no other requirement–no "magic," so to speak.

```bash
$ cat handler.go
```

```go
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
```

## Create templates

Go has powerful template engine out of box.
Now we need to create folder `templates` and put all needed templates inside.

- Create a file named `index.html` inside the `templates` folder and put the following code inside it.

```go
{{ define "content" }}
<div class="container">
    <p> All Posts</p>
</div>
<table class="table">
	<thead class="thead-dark">
	<tr>
		<th>ID</th>
		<th>Name</th>
		<th>View</th>
		<th>Edit</th>
		<th>Delete</th>
	</tr>
	</thead>
	<tbody>
  {{ range . }}
	<tr>
		<tr scope="col">{{ .Id }}</tr>
		<tr scope="col"> {{ .Name }} </tr>
		<tr scope="col"><a href="/show?id={{ .Id }}">View</a></tr>
		<tr scope="col"><a href="/edit?id={{ .Id }}">Edit</a></tr>
		<tr scope="col"><a href="/delete?id={{ .Id }}">Delete</a><td>
	</tr>
  {{ end }}
	</tbody>
</table>
{{ end }}
```

- Next create a file named `new.html` inside the `templates` folder and put the following code inside it.

```go
{{ define "content" }}
<h2>New Blog Post</h2>
<form class="form" method="POST" action="insert">
	<label> Name </label><input type="text" name="title" /><br />
	<input type="submit" value="Save post" />
</form>
{{ end }}
```

- At last, create `edit.html` template file for update blog post, so again create this file in `templates` folder:

```go
{{ define "content" }}
<h2>Edit Plog Post</h2>
<form method="POST" action="update">
	<input type="hidden" name="uid" value="{{ .Id }}" />
	<label> Name </label><input type="text" name="name" value="{{ .Name }}"  /><br />
	<input type="submit" value="Save Blog Post" />
</form><br />
{{ end }}
```

## Few things about static files
TBD Is it need for 20 minutes article?

## Authentication using Okta
To divide our readers and writers of blog posts, it make sense that writers must have account and
authorized before editing posts, but readers can read without any credentials.
Let's add Okta login for blog authors:

```bash
$ cat auth_handler.go
```

```go
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/okta/okta-jwt-verifier-golang"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var (
	state = "ApplicationState"
	nonce = "NonceNotSetYet"
	sessionStore = sessions.NewCookieStore([]byte("okta-hosted-login-session-store"))
)


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	type customData struct {
		Profile         map[string]string
		IsAuthenticated bool
		BaseUrl         string
		ClientId        string
		Issuer          string
		State           string
		Nonce           string
	}

	issuerParts, _ := url.Parse(os.Getenv("ISSUER"))
	baseUrl := issuerParts.Scheme + "://" + issuerParts.Hostname()

	data := customData{
		Profile:         getProfileData(r),
		IsAuthenticated: isAuthenticated(r),
		BaseUrl:         baseUrl,
		ClientId:        os.Getenv("CLIENT_ID"),
		Issuer:          os.Getenv("ISSUER"),
		State:           state,
	}
	tpl := template.Must(template.ParseFiles("templates/layout.html", "templates/login.html"))
	tpl.Execute(w, data)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "okta-hosted-login-session-store")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	delete(session.Values, "id_token")
	delete(session.Values, "access_token")

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	type customData struct {
		Profile         map[string]string
		IsAuthenticated bool
	}

	data := customData{
		Profile:         getProfileData(r),
		IsAuthenticated: isAuthenticated(r),
	}
	tpl := template.Must(template.ParseFiles("templates/layout.html", "templates/profile.html"))
	tpl.Execute(w, data)
}

func AuthCodeCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Check the state that was returned in the query string is the same as the above state
	if r.URL.Query().Get("state") != state {
		fmt.Fprintln(w, "The state was not as expected")
		return
	}
	// Make sure the code was provided
	if r.URL.Query().Get("code") == "" {
		fmt.Fprintln(w, "The code was not returned or is not accessible")
		return
	}

	exchange := exchangeCode(r.URL.Query().Get("code"), r)

	session, err := sessionStore.Get(r, "okta-custom-login-session-store")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, verificationError := verifyToken(exchange.IdToken)

	if verificationError != nil {
		fmt.Println(verificationError)
	}

	if verificationError == nil {
		session.Values["id_token"] = exchange.IdToken
		session.Values["access_token"] = exchange.AccessToken

		session.Save(r, w)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func exchangeCode(code string, r *http.Request) Exchange {
	authHeader := base64.StdEncoding.EncodeToString(
		[]byte(os.Getenv("CLIENT_ID") + ":" + os.Getenv("CLIENT_SECRET")))

	q := r.URL.Query()
	q.Add("grant_type", "authorization_code")
	q.Add("code", code)
	q.Add("redirect_uri", "http://localhost:8081/authorization-code/callback")

	url := os.Getenv("ISSUER") + "/v1/token?" + q.Encode()

	req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", "Basic "+authHeader)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/x-www-form-urlencoded")
	h.Add("Connection", "close")
	h.Add("Content-Length", "0")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var exchange Exchange
	json.Unmarshal(body, &exchange)

	return exchange

}

func isAuthenticated(r *http.Request) bool {
	session, err := sessionStore.Get(r, "okta-custom-login-session-store")

	if err != nil || session.Values["id_token"] == nil || session.Values["id_token"] == "" {
		return false
	}

	return true
}

func getProfileData(r *http.Request) map[string]string {
	m := make(map[string]string)

	session, err := sessionStore.Get(r, "okta-custom-login-session-store")

	if err != nil || session.Values["access_token"] == nil || session.Values["access_token"] == "" {
		return m
	}

	reqUrl := os.Getenv("ISSUER") + "/v1/userinfo"

	req, _ := http.NewRequest("GET", reqUrl, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", "Bearer "+session.Values["access_token"].(string))
	h.Add("Accept", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	json.Unmarshal(body, &m)

	return m
}

func verifyToken(t string) (*jwtverifier.Jwt, error) {
	fmt.Println(nonce)
	tv := map[string]string{}
	tv["nonce"] = nonce
	tv["aud"] = os.Getenv("CLIENT_ID")
	jv := jwtverifier.JwtVerifier{
		Issuer:           os.Getenv("ISSUER"),
		ClaimsToValidate: tv,
	}

	result, err := jv.New().VerifyIdToken(t)

	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if result != nil {
		return result, nil
	}

	return nil, fmt.Errorf("token could not be verified: %s", "")
}

type Exchange struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int    `json:"expires_in,omitempty"`
	Scope            string `json:"scope,omitempty"`
	IdToken          string `json:"id_token,omitempty"`
}
```

First you need to create Okta Developer Account [https://developer.okta.com/signup/](https://developer.okta.com/signup/) and verify you email.
An Okta Application, configured for Web mode. This is done from the Okta Developer Console and you can find instructions [here](https://developer.okta.com/authentication-guide/implementing-authentication/auth-code#1-setting-up-your-application). When following the wizard, use the default properties. They are are designed to work with our sample applications.

After successfully created account, you need to copy `.env.example` to `.env` and replace with your values:

```bash
CLIENT_ID=
CLIENT_SECRET=
ISSUER=https://{yourOktaDomain}.com/oauth2/default
```

Now we need to `login.html` template in `templates` folder:

```go
{{ define "content" }}

<div id="sign-in-widget"></div>
<script type="text/javascript">
	var config = {};
	config.baseUrl = "{{ .BaseUrl }}";
	config.clientId = "{{ .ClientId }}";
	config.redirectUri = "http://localhost:8081/authorization-code/callback";
	config.authParams = {
		{{/*issuer: "{{ .Issuer }}",*/}}
		responseType: 'code',
		{{/*state: "{{ .State }}" || false,*/}}
		display: 'page',
		scope: ['openid', 'profile', 'email'],
		{{/*nonce: '{{ .Nonce }}',*/}}
	};
	new OktaSignIn(config).renderEl(
		{ el: '#sign-in-widget' },
		function (res) {
		}
	);
</script>

{{ end }}
```

## Routers

Usually `routers.go` or just `main.go` place where models, handlers and tepmlates are all together in one synergy:

```go
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
```

## Fin

🎉 🎉 🎉 Congratulations! You build your first blog in go. 

## Do Even More with Go

Check out more tutorials on Go:

- TBD

If you have any questions, please don’t hesitate to leave a comment below, or ask us on our [Okta Developer Forums](https://devforum.okta.com/). Don't forget to follow us on Twitter [@OktaDev](https://twitter.com/oktadev), on [Facebook](https://www.facebook.com/oktadevelopers) and on [YouTube](https://www.youtube.com/channel/UC5AMiWqFVFxF1q9Ya1FuZ_Q/featured)!
