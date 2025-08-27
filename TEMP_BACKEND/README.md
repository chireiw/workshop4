TEMP_BACKEND - Simple Fiber Auth & Points API

Overview

This small Go (Fiber) service provides user registration, authentication (JWT), user profile and a points transfer feature backed by SQLite (GORM).

Run locally

1. cd to the project folder:
   cd dev_project/TEMP_BACKEND
2. Download modules and run:
   go mod tidy
   go run main.go
3. The server listens on :3000

Endpoints

1) POST /register
- Description: Create a new user
- Request JSON:
  {
    "email": "user@example.com",
    "password": "secret",
    "first_name": "First",
    "last_name": "Last",
    "phone": "0123456789",
    "birthday": "2000-01-01"
  }
- Response: { "id": <id>, "email": "...", "points": 100 }

2) POST /login
- Description: Authenticate and receive a JWT
- Request JSON: { "email": "user@example.com", "password": "secret" }
- Response: { "token": "<jwt>" }

3) GET /me
- Description: Return current user info
- Authorization: Bearer <token>
- Response: user object (password omitted)

4) POST /transfer
- Description: Transfer points from authenticated user to another user
- Authorization: Bearer <token>
- Request JSON: { "to_email": "other@example.com", "amount": 10 }
- Response: { "status": "ok" }

Database

- SQLite file: temp_backend.db (created in project root)
- Migration: Auto-migrated by GORM on startup

Swagger / OpenAPI

- OpenAPI JSON is available in ./docs/openapi.json
- The service serves static ./docs at /swagger; to view with Swagger UI, host Swagger UI locally and point it at /swagger/openapi.json, or copy a Swagger UI distribution into ./docs.

Notes

- JWT secret is a hardcoded example; rotate and move to env vars for production
- No email verification or rate limiting included
- Consider adding input validation and tests before production use

Repository

This folder is committed to: https://github.com/chireiw/workshop4/tree/main/TEMP_BACKEND
