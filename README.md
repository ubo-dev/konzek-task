# Konzek-task

This Golang project provides to reach CRUD operations for Tasks as Authorized User with JWtToken.

## Endpoints

```html
POST /login - login
POST /user  - registiration(creating user)
GET  /user  - getting all users

POST /task - creating task
PUT  /task - updating task
GET  /task - getting all tasks
GET  /task/{id} - getting task by id
DELETE /task/{id} - deleting task by id
```


### Requirements And Versions

---
- PostgreSQL
- go 1.21.5
-	github.com/golang-jwt/jwt/v4 v4.5.0
-	github.com/gorilla/mux v1.8.1
-	github.com/lib/pq v1.10.9
-	golang.org/x/crypto v0.17.0
-	github.com/davecgh/go-spew v1.1.1 
-	github.com/pmezard/go-difflib v1.0.0 
-	github.com/stretchr/testify v1.8.4
-	gopkg.in/yaml.v3 v3.0.1

### Run & Build


#### MakeFile
```
build:
	@go build -o bin/konzek-task

run: build
	@./bin/konzek-task

test:
	@go test -v ./..
```
___
*$PORT: 3000*
```ssh
$ cd konzek-task
$ make run
```
