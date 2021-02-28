# login-api

login-api is a library for checking a user's credentials and redirecting them to the client webpage.

Read more about this project [here](https://www.techtrek.io/getting-started-with-go/).

## Requirements
- Install go (instructions [here](https://golang.org/doc/install))

## Running Locally

- Within the project's directory, run `go run main.go`
- Using postman (or equivalent), hit the local url `http://localhost:10000` with the following body: 
```json
{
    "username": {{username}},
    "password": {{password}},
    "token": {{token}}
}
```

## Running Unit Tests
```bash
go test *.go
```
## Usage
Call the api with the following body:

```json
{
    "username": {{username}},
    "password": {{password}},
    "token": {{token}}
}
```

