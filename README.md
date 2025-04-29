# Event REST API

A simple RESTful API for managing events and users, built with [Go](https://golang.org/) and [Gin](https://gin-gonic.com/).

---

## Features

- User registration and login (JWT authentication)
- CRUD operations for events
- Swagger documentation
- Passwords hashed with bcrypt
- Modular project structure

---

## Tech Stack

- [Go](https://golang.org/)
- [Gin Web Framework](https://gin-gonic.com/)
- [Swaggo for Swagger](https://github.com/swaggo/swag)
- [JWT Authentication](https://github.com/golang-jwt/jwt)

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/your-username/event-rest-api.git
cd event-rest-api
```

### 2. Install dependencies
```bash
go mod tidy
```

### 3. Generate swagger docs
```bash
swag init --dir cmd/api --parseDependency --parseInternal --parseDepth 1
```

### 4. Run the application
#### Install air for hot reload
```bash
go install github.com/cosmtrek/air@latest
```
#### Finally run the application
```bash
air
```


