---
layout: blog_post
title: "Go web app in 20 minutes"
author: a_soldatenko
description: "This tutorial explains how to build a web application using Go in less than 20 minutes."
tags: [go, html, blog]
---

# Go web app in 20 minutes

## Go and web applications

Building web apps in Go is simple, fun, and most of all: performant. If you come from an interpreted background (like myself, with Python), you might be surprised as how fast your Go apps can be! If you're all about minimizing your server-side latencies while still having a fun time writing code: you'll probably love Go as much as I do! =)

In this tutorial, I'll walk you through building a simple blog in Go using popular open source libraries (and an API service). If you'd like to see how to craft a simple, real-world Go application, continue reading!

## Prerequisites

Before we dig into building web app in Go, please make sure you have installed the latest version of Go on your operating system. To see which version you have installed you can run:

```bash
$ go version
go version go1.11.2 darwin/amd64
```

If you don't have Go installed, please go [Install the Go tools](https://golang.org/doc/install#install) and make sure you download the [latest stable version](https://golang.org/dl/) for your operating system. In this tutorial I'll be using Go 1.11.2, any version of Go 1.11 or higher should be OK though.

Once you've got Go installed, try running `$ go version` in your terminal again and make sure you've got it working properly.

## Write Your First Go Web App

As I mentioned before, in this tutorial, you'll be building a blog. So the first thing you'll want to do is go create a `go-blog` folder to hold your project source and then create your first go file, `server.go`, in that folder.

```bash
mkdir go-blog
cd go-blog
touch server.go
```

Now, copy the code below and paste it into your new `server.go` file. This file will contain the main source of your web server (app).

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

A *handler* receives and processes HTTP requests sent from the client. In this case, your handler function also uses the Go template engine (which is built into the language) to generate some HTML and finally bundles data into the HTTP response which is eventually sent back to the client.

Simple, right?

Now that your code is ready to roll, you can build and run your first go web app!

```bash
go run server.go
```

**NOTE**: `go run` will both build (compile) and run your go program. It's a great way to test your application while in development, since you don't need to compile and then run your app, you can just do it all at once!

Now, if you open your browser and go to `http://localhost:8081/okta`, you should see:

```
Hello okta!
```

### Sidenote: Go Modules

Since Go 1.11, there is *finally* an official solution for dependency management: [Go modules](https://github.com/golang/go/wiki/Modules)! While this feature is still technically experimental, it's already been widely adopted by the Go community and is considered a best practice.

The rest of this tutorial assumes you're using Go 1.11, primarily because you'll be creating your Go blog as a module, which means you won't need to place your project folder in a special location: you can store your Go code anywhere you want!

## Initialize Your Go Project

Since you've already created a `go-blog` project folder above, let's continue using this folder. All you need to do to initialize it is to run:

```bash
$ go mod init github.com/<your-github-username>/go-blog
```

The `go mod init` command (and be sure to substitute in your proper GitHub username) will initialize this folder as a Go module, meaning all of your future dependency management operations will work properly.

When you run the `go mod init` command, Go will create a new `go.mod` file in your project folder. This file is what Go uses to detect and work with module dependencies.

Now run the following commands to prepare your environment:

```bash
rm server.go
touch main.go
```

Your blog's main function (where everything starts) will be placed in the `main.go` file as we go along. The old `server.go` file you had created before can be safely deleted as we won't need it any longer.

## Create Your Go Database Models

In Go, you will typically use a `struct` for mapping data models to a database. In this blog app, you'll need two models:  `Author` (which represents a blog author), and `Post` (which represents a blog post).

Go ahead and create a `main.go` file and paste in the following code:

```go
package main

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type Author struct {
	Name        string
	Email       string
}

func (u Author) String() string {
	return fmt.Sprintf("Author<%s %s>", u.Name, u.Email)
}

type Post struct {
	Id          int64
	Title       string
	Content     string
	AuthorEmail string
}

func (s Post) String() string {
	return fmt.Sprintf("Post<%d %s %s>", s.Id, s.Title, s.AuthorEmail)
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

func createAuthorProfile(db *pg.DB, profile map[string]string) {
	fmt.Println(profile)
	author := &Author{
		Name: profile["name"],
		Email: profile["email"],
	}
	db.Model(author).Where("email = ?", author.Email).SelectOrInsert()
}
```

You might have noticed that in this file we're using the [go-pg](https://github.com/go-pg/pg) library, which is a popular ORM for working with [Postgres](https://www.postgresql.org/) in Go. I <3 Postgres and hope you do too!

**PS**: If you're looking for a way to scale your existing Postgres databases for very high throughput loads, be sure to check out [Citus Data](https://www.citusdata.com/). They've built a very cool Postgres extension which makes scaling Postgres simple. I highly recommend it.


Now that your database code has been written, why not try to run it? Go to your terminal and try to run the `main.go` program you just built.

```bash
go run main.go
```

If you can see error like this:

```bash
main.go:7:5: cannot find package "github.com/go-pg/pg" in any of:
```

Try to run `go mod init github.com/andriisoldatenko/go-blog` and than run again `go run main.go`. Please don't forget to replace `andriisoldatenko` to you username.
AHA! Did you see what happened there? If all went well, you should have seen output like the following:

```bash
go run main.go
go: finding github.com/go-pg/pg/orm latest
go: downloading github.com/go-pg/pg v6.15.1+incompatible
go: finding github.com/jinzhu/inflection latest
...
```

Even though you didn't explicitly install the `go-pg` library, when Go looked at your code and attempted to build your project, it noticed that the `go-pg` dependency was missing, and automatically downloaded it from GitHub for you! Pretty amazing, right?

Anyhow, you should also have gotten an error, however, that looks something like this:

```bash
panic: dial tcp [::1]:5432: connect: connection refused
```

This is expected since you didn't configure your database at all, and (probably) didn't have a Postgres database running already. Let's get this all set up.

## Set Up a Postgres Database


You have a couple of options on how to run Postgres.

If you're on a Mac, you can always use Postgres via [Postgres app](https://postgresapp.com/). It's the simplest way to get Postgres running on your Mac.

If you're on some sort of *nix OS, you can always install Postgres via your favorite package manager (on Debian-based distributions this means you'd likely install the `postgresql` package).

Finally, if you're super hip and into [containerizing everything](https://www.youtube.com/watch?v=gES4-X6y278), you might want to just install Postgres via [Docker](https://www.docker.com/).

Installing and managing Postgres is a bit out of scope for this article, so I'll let you decide which method to use. But… For simplicity's sake, I'll show you how to set up Postgres in a Docker container (because that's how I roll).

```bash
$ docker run --name postgresql -itd --restart always \
    --publish 5432:5432 \
    --volume /Users/andrii/docker/postgresql:/var/lib/postgresql \
    --env 'DB_USER=blog' --env 'DB_PASS=blog_secret_password' \
    --env 'DB_NAME=blog_db' \
    sameersbn/postgresql:10
```

**NOTE**: You'll notice that in my `docker` command I'm using some specific folders to store my data. I'm also defining a database username and password. You'll want to adjust these settings accordingly for your environment.

Once this command has finished running, you can then check to ensure your shiny new Postgres database is ready to accept incoming connections by running the command below.

```bash
docker logs -f postgresql
# few missing lines
2018-11-21 16:20:04.714 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
2018-11-21 16:20:04.714 UTC [1] LOG:  listening on IPv6 address "::", port 5432
2018-11-21 16:20:04.718 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
2018-11-21 16:20:04.801 UTC [1105] LOG:  database system was shut down at 2018-11-21 16:20:04 UTC
2018-11-21 16:20:04.834 UTC [1] LOG:  database system is ready to accept connections
```

Because I'm paranoid and like to make sure things are working as expected, I also like to connect to my Postgres DB manually using the `psql` CLI tool as well:

```bash
psql -h localhost -U blog -d blog_db -p 5432
Password for user blog:
psql (10.5, server 10.4 (Ubuntu 10.4-2.pgdg18.04+1))
Type "help" for help.

blog_db=> \dt;
Did not find any relations.
blog_db=>
```

As you can see above, I was able to connect to Postgres and verify that there are no tables in the `blog` DB just yet. So… Let's work on fixing that. =)

## Initialize Postgres Tables


Now that you have Postgres running, you need to do initialize your database tables. Once that's done, things should be workable!

In order to get this working, you need to tell your application how to properly connect to your shiny new Postgres database. To do this, open up `main.go` and find your DB connection settings (they should look like this):

```go
func DBConn() (db *pg.DB) {
	db = pg.Connect(&pg.Options{
		Database: "blog_db",
		User: "blog",
		Password: "blog_secret_password",
	})
	return db
}
```

Now, substitute in the appropriate connection settings for your Postgres instance.

If you now run `go run main.go`, you can see that your tables will now have been created to hold your application's authors and posts.

```bash
Author<1 admin >
[Author<1 admin >]
Post<1 Cool story Author<1 admin >>
```

## Create the Go Handlers

Now that the basic database models are ready and functional, let's go ahead and create the handlers. The handlers are the Go functions that will support the blog's functionality: creating posts, viewing posts, etc.

You'll need to create a handler for each piece of functionality in the app.

To get started, copy and paste the code below into a new file named `handlers.go`.

```go
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


You can see by looking through the code above that:

- The `AllPosts` function displays all posts by pulling them out of the database and displaying them in a template (which hasn't yet been built).
- The `NewPost` function displays the HTML form which allows a user to create a new post.
- The `InsertPost` function takes a new post and stores it in the database.
- Finally, the `EditPost` function allows a user to edit a post.

**NOTE**: I didn't implement post deletion here. This is something you should try to implement yourself! It'll be fun! Hint: here's a [documentation link](https://godoc.org/github.com/go-pg/pg#DB.Delete) you may find useful.

## Create Go Templates

One of the things I love most about Go is that it comes with an out-of-the-box template engine that is really powerful.

To make use of Go's templating engine, all you need to do is create a new folder named `templates` in your project directory and define your template files inside.

Firstly, create `templates/layout.html` and paste the following code into this file.

```go
<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=9">
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	<link href="/static/css/bootstrap.min.css" rel="stylesheet">
	<link href="/static/css/font-awesome.min.css" rel="stylesheet">

  	<title>Blog</title>
</head>
<body>
<nav class="navbar navbar-expand-lg navbar-light bg-light">
	<div class="collapse navbar-collapse" id="navbarSupportedContent">
		<ul class="navbar-nav mr-auto">
			<li class="nav-item active">
				<a class="nav-link" href="#">Home <span class="sr-only">(current)</span></a>
			</li>
			<li class="nav-item active">
				<a class="nav-link" href="/new/">Create Post</a>
			</li>
	{{ if .IsAuthenticated }}
			<li class="nav-item">
				<a class="nav-link" href="/logout/">Logout</a>
			</li>
			<li class="nav-item">
				<a class="nav-link" href="/profile/">{{ .Author.name }}</a>
			</li>
	{{ else }}
			<li class="nav-item">
				<a class="nav-link" href="/login/">Login</a>
			</li>
	{{ end }}
		</ul>
	</div>
</nav>
<div class="container">

{{ block "content" . }}{{ end }}

</div> <!-- /container -->

<script src="/static/js/jquery-2.1.1.min.js"></script>
<script src="/static/js/bootstrap.min.js"></script>
</body>
</html>
```

This is a base template that all other templates will inject their HTML into. Keeping it separated like this makes it easier to build complex UIs with minimal code reuse.

Next, create a file named `templates/index.html` and paste the following code inside of it. This file will hold all of the app's home page markup. This page will be used to display all blog posts (ordered by date, descending).

```html
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
  {{ range .Posts }}
	<tr>
		<th scope="col">{{ .Id }}</th>
		<th scope="col"> {{ .Title }} </th>
		<th scope="col"><a href="/show?id={{ .Id }}">View</a></th>
		<th scope="col"><a href="/edit?id={{ .Id }}">Edit</a></th>
		<th scope="col"><a href="/delete?id={{ .Id }}">Delete</a><td>
	</tr>
  {{ end }}
	</tbody>
</table>
{{ end }}
```


Now, create a file named `templates/new.html` and paste in the following code. This markup will be displayed if a user wants to create a new post.

```html
{{ define "content" }}
<div class="container">
	<p>New Blog Post</p>
	<form action="/new/insert/" method="post">
		<div class="form-group">
			<input type="text" name="title" class="form-control" placeholder="Enter title">
		</div>
		<div class="form-group">
			<textarea class="form-control" rows="3"  placeholder="Enter your story"></textarea>
		</div>
		<input type="submit" value="Save post" />
	</form>
</div>
{{ end }}
```

Finally, create a file named `templates/edit.html`. This file will hold the markup that renders the post editing UI. Paste the following code into this file.

```html
{{ define "content" }}
<h2>Edit Plog Post</h2>
<form method="POST" action="update">
	<input type="hidden" name="uid" value="{{ .Id }}" />
	<label> Name </label><input type="text" name="name" value="{{ .Name }}"  /><br />
	<input type="submit" value="Save Blog Post" />
</form><br />
{{ end }}
```

## Add User Authentication to Your Go App

Now that you've got a minimalistic blog with a back-end and front-end, you need to add in the concept of users and user authentication to your app. You've got to add the ability to let users log into your blog and create a post!

There are lots of different ways to do this yourself, but the reality of the situation is that there are a lot of risks that come with implementing user authentication yourself: standards change, it's easy to miss something important that can cause a security risk later, etc. It's a lot safer to just… Not roll your own authentication.

That's why instead of rolling user authentication manually, I'll show you how to accomplish the same thing in a simpler, safer way using Okta's [free API service](https://developer.okta.com) for user management.

To get started, create a new file named `auth_handler.go`. This file will hold your auth-related handler functions.

```go
package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/okta/okta-jwt-verifier-golang"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	state = "ApplicationState"
	nonce = "NonceNotSetYet"
	sessionStore = sessions.NewCookieStore([]byte("okta-hosted-login-session-store"))
)


type Exchange struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int    `json:"expires_in,omitempty"`
	Scope            string `json:"scope,omitempty"`
	IdToken          string `json:"id_token,omitempty"`
}

type CustomData struct {
	Author          map[string]string
	IsAuthenticated bool
	Email           string
}

func GenerateNonce() (string, error) {
	nonceBytes := make([]byte, 32)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("could not generate nonce")
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}

func AuthHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		author := getAuthorData(r)
		data := CustomData{
			Author:         author,
			IsAuthenticated: isAuthenticated(r),
		}
		context.Set(r, "data", data)
		context.Set(r, "email", author["email"])
		next.ServeHTTP(w, r)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	nonce, _ = GenerateNonce()
	var redirectPath string

	q := r.URL.Query()
	q.Add("client_id", os.Getenv("CLIENT_ID"))
	q.Add("response_type", "code")
	q.Add("response_mode", "query")
	q.Add("scope", "openid profile email")
	q.Add("redirect_uri", "http://localhost:8081/authorization-code/callback")
	q.Add("state", state)
	q.Add("nonce", nonce)
	redirectPath = os.Getenv("ISSUER") + "/v1/authorize?" + q.Encode()

	http.Redirect(w, r, redirectPath, http.StatusMovedPermanently)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "okta-hosted-login-session-store")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	delete(session.Values, "id_token")
	delete(session.Values, "access_token")
	context.Clear(r)
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	data := context.Get(r, "data")
	tpl, _ := template.ParseFiles("templates/layout.html", "templates/profile.html")
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

	session, err := sessionStore.Get(r, "okta-hosted-login-session-store")
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

	db := DBConn()
	createAuthorProfile(db, getAuthorData(r))
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
	session, err := sessionStore.Get(r, "okta-hosted-login-session-store")

	if err != nil || session.Values["id_token"] == nil || session.Values["id_token"] == "" {
		return false
	}

	return true
}

func getAuthorData(r *http.Request) map[string]string {
	m := make(map[string]string)

	session, err := sessionStore.Get(r, "okta-hosted-login-session-store")

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
```


The file you just created now holds all of your authentication-related operations, including:

- Login (via OpenID Connect),
- Logout (via OpenID Connect),
- Profile management (modifying your user profile), and
- Helper functions to determine whether or not a user is authenticated or not

The way all this user authentication works is via OpenID Connect's [Authorization Code flow](https://developer.okta.com/authentication-guide/implementing-authentication/auth-code). This flow allows a user to sign into a server-side web application in a secure and standardized way.

If you'd like to see how these implementation details work, look through the code above and try to understand how the handlers are working. Reading code is one of the best ways to familiarize yourself with new concepts.

Back to business though.

Before you can run your new code and test everything out, you'll need to actually set up and configure Okta.

To do this, you need to first go create a free [Okta Developer account](https://developer.okta.com/signup/).

Once you've created your account and are in the Okta dashboard, copy down your **Org URL** value from the top-right hand side of your dashboard page. This will be needed in a moment.

Then visit the **Applications** tab, and click the **Add Application** button. Select the **Web** icon and click **Next**. On the following page, you'll need to define your app-specific settings (these tell Okta what type of app you're building so it knows how to secure it properly).

- **Name**: Set the name to whatever you like. I prefer to name my Okta applications the same name as my project.
- **Base URIs**: Set this to `http://localhost:8081`
- **Login redirect URIs**: Set this to `http://localhost:8081/authorization-code/callback`

Next, click **Done** and your new application will be created. Copy down the **Client ID** and **Client Secret** values on this page. You'll need these in a moment.

Now what you need to do is create a file named `.env` in your project directory. This file will hold all of your super-secret application values: like your Okta application's ID and secret, etc. Create this file and insert the following values into it.

```bash
export CLIENT_ID=<your client ID>
export CLIENT_SECRET=<your client secet>
export ISSUER=<your org URL>/oauth2/default
```

To "activate" these environment variables so your app can see them, run the command below:

```bash
source .env
```


Now, you need to create a `templates/login.html` file to hold the login page markup. Paste the following code inside this new file.

```html
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

## Create Your Go Router

So far you've built out the blog functionality in templates, handlers, and models. But you still need to plug it all together with a router. A router simply maps URL requests from a user to handler code.

Open up `main.go` and insert the following code (replacing what was there before).

```go
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
```

## Finalizing Your Go Application

At this point, you've successfully built your first Go web application! Congratulations!

To test it out, run `go run *.go`, then visit `http://localhost:8081` in your browser. If you everything is working properly you should now be able to use your new blog!

**NOTE**: `go run *.go` means to run all `*.go` files in current folder, we need this because we have one `main` package inside different go files.

If you liked this small Go tutorial, you might want to check out some of our other articles below (or [follow us on Twitter](https://twitter.com/oktadev), [Facebook](https://www.facebook.com/oktadevelopers), and [YouTube](https://www.youtube.com/channel/UC5AMiWqFVFxF1q9Ya1FuZ_Q/featured)). We publish all sorts of technical tutorials and guides you may find interesting.

- [Build a Single-Page App with Go and Vue](https://developer.okta.com/blog/2018/10/23/build-a-single-page-app-with-go-and-vue)
- [OpenID Connect - A Primer](https://developer.okta.com/blog/2017/07/25/oidc-primer-part-1)
- [Learn JavaScript in 2019](https://developer.okta.com/blog/2018/12/19/learn-javascript-in-2019)
- [Build and Test a React Native App with TypeScript and OAuth 2.0](https://developer.okta.com/blog/2018/11/29/build-test-react-native-typescript-oauth2)
- [Build a Desktop App with Electron and Authentication](https://developer.okta.com/blog/2018/09/17/desktop-app-electron-authentication)

If you have any questions, please don’t hesitate to leave a comment below, or ask us on our [Okta Developer Forums](https://devforum.okta.com/).