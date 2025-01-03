# Golang Little Twitter
## Description

This is an project which implements end-points for a little mimic twitter. Aplication that contains enpoints according to uala backend chanllege.
It exposes endpoints to post a new tweets, set a new follower and get timeline for a user.


## Overview
This application was designed following clean architecture and uses the principles of Robert Martin. It contains endpoints for:
- Save a new tweet
- Add new user follower
- Get the timeline for a specific user


Rule of Clean Architecture by Uncle Bob

- Independent of Frameworks.
- Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element.
- Independent of UI. The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.
- Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.
- Independent of any external agency. In fact your business rules simply don’t know anything at all about the outside world.

More at https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html

This project has 4 Domain layer :

- Entity Layer
- Repository Layer
- Usecase/Service Layer
- Delivery Layer


## Content
- [Internal Architecture](#internal-architecture)
- [Project structure](#project-structure)
- [Dependency Injection](#dependency-injection)
- [Running Application](#running-application)
- [Swagger Documentation](#swagger-documentation)

## Internal Architecture]
#### The diagram:

![l-twitter architecture](https://github.com/ibanezv/l_twitter_api/raw/main/docs/img/ltwitter.arch.drawio.png)

#### The diagram in detail:

![l-twitter detailed architecture](https://github.com/ibanezv/l_twitter_api/raw/main/docs/img/ltwitter.arch.detail.png)

### How To Run This Project

> Make Sure you have run the ltwitter.sql in your posgressql


## Project structure
### `cmd/api/main.go`
Configuration and logger initialization. Then the main function continues in
`internal/app/app.go`.

### `config`
Configuration. First, `config.yml` is read, then environment variables overwrite the yaml config if they match.
The config structure is in the `config.go`.
The `env-required: true` tag obliges you to specify a value (either in yaml, or in environment variables).
It is assumed that default values are in yaml, and security-sensitive variables are defined in ENV.

### `docs`
Swagger documentation. Auto-generated by [swag](https://github.com/swaggo/swag) library.

### `internal/app`
There is always one _Run_ function in the `app.go` file.

This is where all the main objects are created.
Dependency injection occurs through the "New ..." constructors (see Dependency Injection).
This technique allows us to layer the application using the [Dependency Injection](#dependency-injection) principle.
This makes the business logic independent from other layers.

The `migrate.go` file is used for database auto migrations.
It is included if an argument with the _migrate_ tag is specified.
For example:

```sh
$ go run -tags migrate ./cmd/api
```

### `internal/app/handlers`
Server handler layer. It contains controllers for :
- Followers
- Timelines
- Twitter

Server routers are written in the same style:
- Handlers are grouped by area of application
- For each group, its own router structure is created, the methods of which process paths
- The structure of the business logic is injected into the router structure, which will be called by the handlers

In `router.go` and above the handler methods, there are comments for generating swagger documentation using [swag](https://github.com/swaggo/swag).

### `models`
Entities of business logic (models) can be used in any layer.

### Use cases /  Service
#### `internal/follower`
#### `internal/timeline`
#### `internal/tweet`
Business logic.
- Methods are grouped by area of application
- Each group has its own structure
- One file - one structure

Repositories are injected into business logic structures

### `repository/repo`
A repository is an abstract storage that business logic works with.

### `pkg/db`
A database abstract (DBStorage)

### `pkg/cache`
A cache engine abstract (Cache)

### `pkg/postgres`
It is a postgres implemenation of DBStorage

### `pkg/redis`
It is a redis implemenation of DBStorage

## Dependency Injection
In order to remove the dependence of business logic on external packages, dependency injection is used.

## Running Application

#### Run tests

```bash
$ make tests
```

### Run the application locally

Steps in order to run code repository

```bash
#install postgres and initialize PostgreSQL
$ sudo mkdir -p /etc/paths.d && echo /Applications/Postgres.app/Contents/Versions/latest/bin | sudo tee /etc/paths.d/postgresapp

#install redis docker image
$ docker pull redis

#run redis docker image
$ docker run --name redis -d redis 

#move to directory
$ cd workspace

# Clone into your workspace
$ git clone https://github.com/ibanezv/l_twitter_api.git

#move to project directory
$ cd l_twitter_api

# copy the example.env to .env
$ cp example.env .env

# Migrate database
$ make migrate-up

# Run the application
$ go run cmd/api/main.go

# Execute the call on terminal
-
$ curl -X 'GET' 'http://localhost:8080/health'
-
$ curl -X 'GET' 'http://localhost:8080/v1/api-twitter/timeline/1' \-H 'accept: application/json'
-
$ curl -X 'POST' 'http://localhost:8080/v1/api-twitter/follow' \
--header 'Content-Type: application/json' \
--data '{
    "user_id":200,
    "user_id_followed":100 }'
-    
$ curl -X 'POST' 'http://localhost:8080/v1/api-twitter/tweet' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": 100,
    "text": "message of user 100",
    "date_time": "2024-12-30T00:00:00Z" }'
```

Running using docker compose

```bash
# Postgres, Redis
$ make compose-up

# Run app with migrations
$ make run
```