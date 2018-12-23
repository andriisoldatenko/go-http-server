---
layout: blog_post
title: "Build an HTTP Server with Go"
author: a_soldatenko
description: "This post walk trough building a simple API service in Go"
tags: [go, gin, rest, jwt, okta]
---

# Build an HTTP Server with Go


## Film application

Nowadays it's very popular to build web applications using JavaScript and HTML and some backend for instance using MVC frameworks like Ruby on Rails of Django. But some times you want to build something small, usefull and secure. Let's imagine one day you decided to create API to share information about films, and you would like to have two kinds of application:
- some endpoints can be protected, and user can register to see it
- another endpoint can be public :)

## Few words about RESTful (HTTP)

Representational State Transfer (REST) is a software approach that declarates a list of rules to be used for building you service. In another words service which apply REST style, allows to access and edit some amount of data by using set of defined operations.

## How it looks like when it's done

In this tutorial I'm going to use my favorite programming language `Go`. app I'll demonstrate how to build simple API using `gin` http server framework and JWT and Okta autorization.

GET 'http://locahost:8080/api/v1/films'
GET 'http://locahost:8080/api/v1/film/<id>/'
GET 'http://locahost:8080/api/v1/film/<id>/'
using JWT
GET 'http://locahost:8080/api/v1/protected/films/'
GET 'http://locahost:8080/api/v1/protected/film/<id>/'

## Database

Usually database refers to data and the way how its organized inside. I prefer to use `PostgreSQL` for all my projects, because it's open source and has ton's of feature. More details we can find in officila documentation. I'll brefly show you how to getting started with PostgreSQL and Docker.
Docker de-facto new standart similar to `git` or another tools.
?? TODO more info about Docker?
Make sure Docker is installed in you environment. https://docs.docker.com/install/

```bash
$ docker --version
Docker version 18.09.0, build 4d60db4
```

Let's go to postgres [official docker hub](https://hub.docker.com/_/postgres/) and pull latest official PostgeSQL docker image:

```bash
$ docker run --name go-server-postgres -e POSTGRES_PASSWORD=postgres -d postgres
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

## Models

Usually in Golang we create `struct`s to define models or tables in databases

```go
type Film struct {
	gorm.Model
	Title        string `gorm:"unique_index"`
	Year         string
	Plot         string `gorm:"size:2048"`
}
```


## Create our first view 

```
```

## What about tests?



## How to debug if something goes wrong?

