# Architecture — C4 Summary for TEMP_BACKEND

This document summarizes the architecture of the TEMP_BACKEND project using C4-style diagrams (PlantUML C4). The diagrams describe system context, containers, and components derived from the current codebase (`main.go`).

Render instructions

- These diagrams use the C4-PlantUML library. To render them:
  - Use the PlantUML CLI with an internet connection (the diagrams include the C4-PlantUML includes), or
  - Use an online PlantUML server/renderer that supports remote !include, or
  - Download the C4-PlantUML set locally and update the `!include` paths.

System Context (Level 1)

```plantuml
@startuml SystemContext
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Context.puml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4.puml

Person(user, "End user", "Uses web UI or API clients")
System_Boundary(s_temp, "TEMP_BACKEND") {
  Container(web_ui, "Web UI (optional)", "React (Vite)", "Browser-based front-end (not included in this repo)")
  Container(api, "API Server", "Go + Fiber", "Provides REST API: /register, /login, /me, /transfer")
  ContainerDb(db, "SQLite DB (temp_backend.db)", "SQLite", "Persists users and points")
}

Rel(user, web_ui, "Uses")
Rel(user, api, "Calls API (Bearer token) or via Web UI")
Rel(web_ui, api, "Calls REST API")
Rel(api, db, "Reads/Writes using GORM")

@enduml
```

Container (Level 2) — API Server decomposition

```plantuml
@startuml Containers
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Container.puml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4.puml

Person(user, "End user", "Authenticated user with JWT")
System_Ext(swagger, "Swagger UI / OpenAPI", "Docs", "Optional API documentation viewer served from ./docs")

Container(api, "API Server", "Go + Fiber", "Implements auth and points transfer")
ContainerDb(db, "SQLite DB (temp_backend.db)", "SQLite", "Stores users and points")

Container_Boundary(apiBoundary, api, "API Server") {
  Component(authHandler, "Auth Handler", "handles /register and /login", "Hashes passwords, creates users, issues JWTs")
  Component(userHandler, "User Handler", "handles /me", "Returns current user data")
  Component(transferHandler, "Transfer Handler", "handles /transfer", "Performs transactional points transfer using GORM transaction")
  Component(tokenService, "Token Service", "JWT generation/validation", "Creates signed JWTs using secret")
  Component(repo, "Repository / DB Layer", "GORM + SQLite", "Data persistence and queries")
}

Rel(user, api, "Uses REST API (Bearer token)")
Rel(api, db, "Reads/Writes via GORM")
Rel(api, swagger, "Serves OpenAPI JSON and static docs")
Rel(authHandler, tokenService, "Generates JWTs")
Rel(transferHandler, repo, "Reads/Writes users within transaction")
Rel(userHandler, repo, "Reads user record")

@enduml
```

Component (Level 3) — Key components and responsibilities

- Auth Handler
  - Parses /register and /login requests
  - Hashes passwords using bcrypt
  - Creates User records via GORM
  - Calls Token Service to produce JWT
- Token Service
  - Produces HS256 JWTs with `sub` claim set to user ID
  - Token lifetime: 72 hours (in code)
  - Uses `jwtSecret` variable (move to env var for production)
- Transfer Handler
  - Parses transfer requests
  - Validates amount and recipient existence
  - Uses GORM Transaction to debit/credit user points atomically
- Repository / DB Layer
  - Single `users` table (GORM auto-migrate)
  - Fields: id, email, password, first_name, last_name, phone, birthday, points, created_at

Operational notes & recommendations

- Move `jwtSecret` into environment variables and provide an example `.env` in the repo.
- Add input validation (email formats, password rules) and better error responses.
- For production, replace SQLite with a server DB (Postgres/MySQL) and add migrations.
- Add tests for transfer concurrency to ensure transactional integrity under load.

References

- C4 Model examples: https://c4model.com/diagrams/
- C4-PlantUML: https://github.com/plantuml-stdlib/C4-PlantUML
