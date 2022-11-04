# Codebox - Go Snippets Sharing App
Framework-less server side CRUD application using built-in Go's `html/template` engine, allowing you to share snippets of code (akin to pastebin.com). 

## Features:
- Creating a shareable code snippet with automatic expiration period
- Sharing a snippet to either (personal) private or public repository
- Basic session-based authentication
- Browsing public snippets
- Resiliency against most common http security concerns (xss, csrf, sql injection)
- Static files are embedded within application using Go's `embed` package

## Running:
After configuring PostgreSQL database locally, navigate to project's root folder and run:
```bash
$ go run ./cmd/web -url="postgres://[username]@localhost:5432/[database_name]"
```
Database snippet can be found at this repo's `DATABASE.txt` file.

## Todo:
- Password reset feature
- Private notes URLs redirection for unauthorized access
- Better snippet management (profile/account route)
- Better syntax highlighting for shared code

## Stack:
- Go 1.18 + `justinas/alice` + `valyala/fasthttp` + `alexedwards/scs` + `jackx/pgx`
- PostgreSQL 14.5

## Credits:
Application was heavily inspired by [Alex Edward's - Let's Go](https://lets-go.alexedwards.net/) book.