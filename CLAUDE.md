# CLAUDE.md

## Project Overview

A personal Go REST API server, built as both a learning project and a real tool. It backs a full micro-frontend (MFE) shell architecture — a single orchestrating shell on port 3000 with individual MFEs on incrementing ports (3001, 3002, ...). Currently the shell and a job-search MFE are in place; planned future MFEs include vehicle maintenance, budgeting, project management, and recipes.

The repo itself is a learning project. Expect rough edges, evolving conventions, and deliberate simplicity. Suggestions are welcome, but the user drives priorities.

## Tech Stack

- **Language:** Go 1.24
- **Router:** `go-chi/chi/v5`
- **Database:** MongoDB (`mongo-driver`)
- **Auth:** JWT (access + refresh tokens via `golang-jwt/jwt`)
- **Password hashing:** argon2id (`golang.org/x/crypto`)
- **Config:** `.env` via `godotenv`
- **UUIDs:** `google/uuid`

## Running Locally

```bash
go run main.go
```

Requires a `.env` file in the repo root with MongoDB connection info and key material (private/public key PEM, refresh token pepper, DB name). The server listens on `localhost:8080`.

No Docker, no test suite yet.

## Architecture

### Layer Pattern

Every domain follows the same three-file structure (use `thing/` as the reference template):

```
<domain>/
  <domain>_handler.go   # HTTP handlers — parse input, call store, write response
  <domain>_routes.go    # Route registration — returns a chi.Router
  <domain>_store.go     # MongoDB data access — CRUD methods
  <domain>_model.go     # Structs (model, request DTOs, etc.)
```

`token/` adds a service layer (`token_service.go`) for business logic (JWT minting, validation, rotation). Other domains go directly handler → store.

### Adding a New Service

Follow this checklist (documented in `README.md`):

1. Create a `<domain>/` package with handler, routes, store, and model files
2. Add a store field to `app.Application` (`app/app.go`)
3. Create the MongoDB collection and store instance in `main.go`
4. Register routes in `server/routes.go`

Use `thing/` as the template. Use `job/` or `recruiter/` as reference for auth-protected services that need the user ID from context.

### Auth

- Access tokens are JWTs validated by `token.RequireAccess` middleware
- Refresh token rotation: old token is revoked, new pair is issued
- User ID is injected into context by middleware; handlers retrieve it via `token.UserIDFromContext(r.Context())`
- Refresh tokens are stored in MongoDB with a TTL index for auto-expiry; access tokens are stateless JWTs (not stored)
- Frontend stores tokens in localStorage; backend communicates via CORS

Protected routes are mounted with `r.With(token.RequireAccess(...)).Mount(...)` in `server/routes.go`.

## Current Conventions

### Naming

- Types: `JobStore`, `ThingHandler`, `RecruiterModel`
- Handler methods: verb + entity — `GetAllJobs`, `CreateJob`, `DeleteThing`
- Store methods: `GetAll`, `GetOne`, `Create`, `Update`, `Delete`
- Route functions: `<Domain>Routes(...)` returning `chi.Router`
- Receiver names: `h` for handlers, `s` for stores

### Responses

- `200 OK` — success with JSON body
- `201 Created` — successful create with JSON body
- `204 No Content` — successful delete or revoke (no body)
- `400 Bad Request` — malformed input
- `401 Unauthorized` — missing or invalid token
- `404 Not Found` — resource not found
- `500 Internal Server Error` — store/unexpected failure

Always set `Content-Type: application/json` before calling `json.NewEncoder(w).Encode(...)`. Log encode failures with `log.Printf`.

### Error Handling

Errors are returned from stores and checked in handlers. Use `fmt.Errorf("context: %w", err)` for wrapping. `log.Printf` for server-side failures before returning HTTP 500.

## Known Inconsistencies (track these)

- `dashboard/` is stubbed out (empty package files) pending a rewrite — don't touch it.

## Budgeting Domain Conventions

- **Monetary amounts** are stored as integers (cents) throughout — `Account.Balance`, `Transaction.Amount`, etc. Frontend handles display formatting. Never use `float64` for money.
- **`Income`** is a bool on `Transaction`, not a category. This is intentional — income has its own dedicated UI section and shouldn't be mixed into expense category lists.
- **`Owner`** is a string (`"mine"` or `"hers"`) on `Account`, `Bill`, and `Transaction`. It replaced the old `Member` entity, which was deleted.
- **`BillMonth`** format is `"2006-01"` (YYYY-MM). It is validated server-side in `PayBill`. Frontend is responsible for passing the correct format.
- **`member/`** package was deleted — do not recreate it.

## What NOT to Do

- Don't add a SQL database without the user explicitly deciding to (MongoDB is the current store; SQL may be added later as a learning exercise).
- Don't add a dependency injection framework (e.g. `wire`, `fx`) — `main.go` already does DI manually and idiomatically by passing stores/services as constructor arguments. A framework adds complexity without solving a real problem at this scale. Worth exploring later as a learning exercise once the basics are solid.
- Don't add tests unless asked — the user is intentionally deferring this.
- Don't restructure the handler/store layering — the current handler→store pattern (with a service layer only where there's real logic, like `token`) is intentional. More elaborate patterns exist in Go: store interfaces for decoupling/mocking, clean/hexagonal architecture with additional layers, etc. These aren't wrong, but they add indirection without a concrete payoff at this scale. Introduce them when there's a real reason (e.g. wanting to mock stores for tests), not speculatively.
- Don't over-engineer. This is a learning project; keep solutions focused and minimal.
- Don't leak implementation details in HTTP error responses (avoid `http.Error(w, err.Error(), ...)` for internal errors — use a generic message and log the real error).

## Notes

- `dateApplied` on jobs is a string by design — all date formatting is handled on the frontend. Don't convert it to `time.Time`.
- The `thing` package is the canonical structural template, but `job`/`recruiter` are better references for auth-aware services.
- The user is learning Go — favor clear, idiomatic Go over clever solutions. Flag non-obvious patterns with a brief explanation.
