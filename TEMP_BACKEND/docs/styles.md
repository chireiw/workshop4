TEMP_BACKEND — Code & Docs Style Guidelines

Purpose

This document describes recommended style and conventions for code and documentation in the TEMP_BACKEND project. Follow these to keep the project consistent and maintainable.

Go code style

- Formatting
  - Run `gofmt` / `go fmt` on all files before commit.
  - Run `go vet` to catch common mistakes.
  - Use a linter (recommended: `golangci-lint`) in CI.

- Project layout
  - Keep a small `main.go` for app wiring. Move handlers and helpers into packages when the project grows (e.g., `auth`, `handlers`, `store`).
  - Use `initDB()` for DB setup and `AutoMigrate` only in local/dev; use real migrations (golang-migrate) for production.

- Naming & structure
  - Export only what is required. Private helpers start with lowercase.
  - Use clear struct tags for JSON and GORM. Example: `Password string `json:"-"`` to avoid leaking secrets in API responses.

- Secrets & config
  - Do not hardcode secrets. Read `JWT_SECRET`, `PORT`, DB config from environment variables.
  - Add `.env.example` with keys and document usage in README.

- Error handling
  - Return proper HTTP status codes (400, 401, 404, 500) with helpful error messages.
  - Avoid leaking internal errors to clients; log details on server side.

- Concurrency & transactions
  - Use DB transactions for multi-row updates (already used in `/transfer`).
  - Add tests to cover concurrent transfers and race conditions.

Documentation style (docs/)

- README.md
  - Keep run instructions, endpoints, and quick examples minimal and accurate.

- OpenAPI
  - Keep `docs/openapi.json` updated when endpoints change. Prefer generating from code if the project grows.

- Diagrams
  - Use Mermaid for diagrams (`docs/*.md` include mermaid blocks). Keep diagrams small and split into files if large.
  - If adding rendered images, place them in `docs/diagrams/` and reference in markdown.

- Architecture & Database docs
  - `architecture.md` should document system context, containers, and components with C4 (Mermaid C4 syntax used here).
  - `database.md` should reflect models and important constraints (unique indexes, FK relations).

Git & commits

- Commit messages
  - Use meaningful, short subject line: `Add transfer endpoint`, `Update README`.

- Branching
  - Use feature branches (feature/...) and pull requests for review.

- .gitignore
  - Exclude artifacts: `temp_backend.db`, `node_modules`, binary builds, `.env`.

Testing & CI

- Add basic unit tests for handlers and a concurrency test for `/transfer`.
- Add CI to run `go test`, linters, and optionally render diagrams.

License & Security

- Add LICENSE if project intended to be shared.
- Rotate secrets, enable secure storage for production, and add rate limiting if exposed publicly.

Contact

- Maintainers: update README with responsible contact or GitHub team.
