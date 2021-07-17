<div align="center">
  <h1 align="center">Go-Auth-Redis-Postgres-API</h1>
  <p align="center">
    <img src="https://img.shields.io/github/languages/code-size/BravoAndres/Go-Auth-Redis-Postgres-API?style=for-the-badge">
    <img src="https://img.shields.io/github/license/BravoAndres/Go-Auth-Redis-Postgres-API?style=for-the-badge">
    <img src="https://img.shields.io/github/go-mod/go-version/BravoAndres/Go-Auth-Redis-Postgres-API?style=for-the-badge">
  </p>
</div>
<br>

`Go-Auth-Redis-Postgres-API` is a RESTful Authentication API built with Go, Fiber framework JWT, Postgres, Redis &amp; more!
This project provides Rest API and was made for fun/practicing to learn about Fiber Framework. Feel free to use/modify or send a pull request.

Contents
========

 * [Features](#features)
 * [Requirements](#requirements)
 * [Usage](#usage)
 * [Todo](#todo)
 * [Built With](#built-with)
 * [License](#license)

### Features
---
+ CRUD Functionality
+ JSON payload request/responses
+ Custom JWT middleware (highly inspired in [fiber middleware](https://github.com/gofiber/jwt))
+ PostgreSQL & Redis for persistance/cache storage
+ UUID
+ Testing
+ Docker (Coming soon! 🚀)
+ Clean Code

### Requirements
---
+ Go 1.16+.
+ PostgreSQL & Redis installed on your system.
+ Copy `.env.example` and rename it to `.env` & set your enviroment variables per usual.

### Usage
---
1. Download and install dependencies with `go mod`
    + `$ go mod download`

2. Running the project (local).
    + `$ go run github.com/BravoAndres/Go-Auth-Redis-Postgres-API/cmd/`

### Todo
---
- [ ] Add Gorm
- [ ] Implement custom zap logger
- [ ] Migrations
- [ ] API Docs (Swagger)
- [ ] Add Tests
- [ ] CI/CD
- [ ] Docker
- [ ] Add simple frontend with Vue.js

### Built With
---
+ [Fiber Framework](https://gofiber.io/)
+ [Zap Logger](https://github.com/uber-go/zap)
+ [JWT Authentication](https://jwt.io/)
+ [PostgreSQL](https://www.postgresql.org/)
+ [Redis](https://redis.io/)
+ [Go](https://golang.org/)

### License
---
[MIT License.](https://tldrlegal.com/license/mit-license)
