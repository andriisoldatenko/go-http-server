---
layout: blog_post
title: "Build an HTTP Server with Go"
author: a_soldatenko
description: "This post walk trough building a simple API service in Go"
tags: [go, gin, rest, jwt, okta, postgresql, models, serializers]
---

# Build an HTTP Server with Go

## Film application

Nowadays it's very popular to build web applications using JavaScript and HTML and some backend for instance using MVC frameworks like Ruby on Rails of Django. But sometimes you want to build something small, useful and secure. Let's imagine one day you decided to create API to share information about films, and you would like to have two kinds of application:
- some endpoints can be protected, and user can register to see it
- another endpoint can be public :)

## Few words about RESTful (HTTP)

Representational State Transfer (REST) is a software approach that declares a list of rules to be used for building you service. In another words service which apply REST style, allows to access and edit some amount of data by using set of defined operations.

## How it looks like when it's done

In this tutorial I'm going to use my favorite programming language `Go`. app I'll demonstrate how to build simple API using `gin` http server framework and JWT and Okta authorization.

GET 'http://localhost:8080/api/films'

## Database

Usually database refers to data and the way how its organized inside. I prefer to use `PostgreSQL` for all my projects, because it's open source and has ton's of feature. More details we can find in official documentation. I'll briefly show you how to getting started with PostgreSQL and Docker.
Docker de-facto new standart similar to `git` or another tools.
?? TODO more info about Docker?
Make sure Docker is installed in you environment. https://docs.docker.com/install/

```bash
$ docker --version
Docker version 18.09.0, build 4d60db4
```

Let's go to postgres [official docker hub](https://hub.docker.com/_/postgres/) and pull latest official PostgreSQL docker image:

```bash
$ docker run --name go-server-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres
$ docker run -it --rm --link go-server-postgres:postgres postgres psql -h postgres -U postgres
Password for user postgres:
psql (11.1 (Debian 11.1-1.pgdg90+1))
Type "help" for help.

postgres=#
```

## Create database and user

```
postgres=# CREATE USER film;
CREATE DATABASE film;
GRANT ALL PRIVILEGES ON DATABASE film TO film;
ALTER USER film WITH PASSWORD 'film';
```

Now you we are ready to use our new database instance.

## Go

## Go dependency management

Since mid of 2018 and go1.11 has been released we can use `go mod` or go modules.

```
$ mkdir go-http-server
$ cd go-http-server
$ go mod init github.com/<your-github-name>/go-http-server
$ ls -la
go.mod    go.sum
```

## Storage

First of all we need to create integration between database and our web application. I prefer to do it using `storage` package. It's very easy, we need to create `storage` folder and create file `db.go`:

```go
package storage

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct.
func Init() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=film dbname=film password=film sslmode=disable")
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)
	DB = db
	return DB
}


func GetDB() *gorm.DB {
	return DB
}
```

## Models

Usually in Golang we create `struct`s to define models or tables in databases. In `go` during web development, we usually using `GORM` it's ORM for Golang.
Let's define our `film` model first and I'll describe what each line does. First create folder `films` and create file `models.go` inside:

```bash
$ ls films
models.go
```

And now copy and paste next snippet to `models.go` file:

```go
package films

import (
	"github.com/andriisoldatenko/go-http-server/storage"
	"github.com/jinzhu/gorm"
)

type FilmModel struct {
	gorm.Model
	Title        string `gorm:"unique_index"`
	Year         string
	Plot         string `gorm:"size:2048"`
}

func FindOneFilm() (FilmModel, error) {
	db := storage.GetDB()
	var model FilmModel
	tx := db.Begin()
	tx.First(&model, 1)
	err := tx.Commit().Error
	return model, err
}
```


Note: `gorm:"unique_index"` is special struct syntax in Golang. Fore more details please visit GORM docs -> [Declaring Models](http://gorm.io/docs/models.html)

Here you can see `FilmModel` struct with all related to film fields. `FindOneFilm` is model manager which do SQL query to database using ORM-like query language and map data to `FilmModel` struct back.

## Serializers

Serializers usually help you to define you response structure and validate if needed before we save data to database.

```go
package films

import "github.com/gin-gonic/gin"

type FilmSerializer struct {
	C *gin.Context
	FilmModel
}

type FilmResponse struct {
	ID             uint                  `json:"-"`
	Title          string                `json:"title"`
}

func (s *FilmSerializer) Response() FilmResponse {
	response := FilmResponse{
		ID:          s.ID,
		Title:       s.Title,
	}
	return response
}
```


## Create views

```
package films

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FilmAnonymousRegister(router *gin.RouterGroup) {
	router.GET("/", FilmList)
}


func FilmList(c *gin.Context) {
	filmModel, err := FindOneFilm()
	if err != nil {
		c.JSON(http.StatusNotFound, errors.New("invalid param"))
		return
	}
	serializer := FilmSerializer{c, filmModel}
	c.JSON(http.StatusOK, gin.H{"films": serializer.Response()})
}
```

## Let's put all together

If we try to build and run `main.go` you can see something like it:

```bash
go run main.go
[GIN-debug] [WARNING] Now Gin requires Go 1.6 or later and Go 1.7 will be required soon.

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /api/films/               --> github.com/andriisoldatenko/go-http-server/films.FilmList (3 handlers)
```


## JSON Web Token Authentication

Film API server authenticates requests using access tokens, which are obtained via the Okta Implicit Flow. In general this flow can be described as follow:
- You need to login using Okta sign-in page and get access and ID tokens as a URI params.
- Than any http client(includes your frontend e.g. AngularJS single page application) can use it.

Now we need to set up your Application choosing `Add Application` inside Okta dev page and create new Application as `SPA`, than click `Done`:
Screenshot 1

To make things simpler I've added few routes to login using Okta Hosted login. IF you open `films/routers.go` and add `LoginHandler` to the bottow of file:

```
func LoginHandler(c *gin.Context) {
	nonce, _ := oktaUtils.GenerateNonce()

	q := c.Request.URL.Query()
	q.Add("client_id", os.Getenv("CLIENT_ID"))
	q.Add("response_type", "token")
	q.Add("scope", "openid")
	q.Add("redirect_uri", "http://localhost:8080/implicit/callback")
	q.Add("state", "ApplicationState")
	q.Add("nonce", nonce)

	redirectPath := os.Getenv("ISSUER") + "/v1/authorize?" + q.Encode()
	c.Redirect(http.StatusMovedPermanently, redirectPath)
}
```

and than register handler using `/login` route in `main.go`:

```
@@ -21,10 +20,12 @@ func main() {
        defer db.Close()

        r := gin.Default()
+       r.GET("/login", films.LoginHandler)
+       r.GET("/implicit/callback", films.AuthCodeCallbackHandler)

```

It helps you to login using you favorite browser and just run api server:

```
$ go run main.go
```

Than navigate to `http://localhost:8080/login/`:
Screenshot 2

After successfully logged in using Okta account you will see redirect:
Screenshot 3

And now we need to store somewhere you JWT access token to test our protected API in next steps:

```
http://localhost:8080/implicit/callback#access_token=<your long JWT token from Okta>&token_type=Bearer&expires_in=3600&scope=openid&state=ApplicationState
```


## Authentication middleware

Now let's add folder `auth` and create file `middlewares.go`:

```
package auth

import (
	"github.com/gin-gonic/gin"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"strings"
)

func isAuthenticated(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return false
	}
	tokenParts := strings.Split(authHeader, "Bearer ")
	bearerToken := tokenParts[1]

	tv := map[string]string{}
	tv["aud"] = "api://default"
	tv["cid"] = os.Getenv("SPA_CLIENT_ID")
	jv := verifier.JwtVerifier{
		Issuer:           os.Getenv("ISSUER"),
		ClaimsToValidate: tv,
	}

	_, err := jv.New().VerifyAccessToken(bearerToken)

	if err != nil {
		return false
	}

	return true
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isAuthenticated(c.Request) {
			err := errors.New("auth error")
			c.AbortWithError(401, err)
		}
	}
}
```

and also we need to apply middleware in `main.go` file and function `main`:

```go
package main

import (
	"github.com/andriisoldatenko/go-http-server/auth"
	"github.com/andriisoldatenko/go-http-server/storage"
	"github.com/gin-gonic/gin"
	"log"
	"os"

	"github.com/andriisoldatenko/go-http-server/films"
	"github.com/jinzhu/gorm"
)


func Migrate(db *gorm.DB) {
	db.AutoMigrate(&films.FilmModel{})
}


func main() {
	db := storage.Init()
	Migrate(db)
	defer db.Close()

	r := gin.Default()
	r.GET("/login", films.LoginHandler)
	r.GET("/implicit/callback", films.AuthCodeCallbackHandler)

	v1 := r.Group("/api")
	v1.Use(auth.AuthMiddleware())
	films.FilmAnonymousRegister(v1.Group("/films"))

	r.Run()
}
```

:tada we are ready to start our new api film server:

```
$ go run main.go
```

```
$ http -v GET http://localhost:8080/api/films/ Authorization:'Bearer <your jwt access key from Okta>'

HTTP/1.1 200 OK
Content-Length: 38
Content-Type: application/json; charset=utf-8
Date: Mon, 24 Dec 2018 13:00:53 GMT

{
    "films": {
        "title": "AAAAAAAAAAAAAAAA"
    }
}
```
