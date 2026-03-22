# Go REST API Server

A REST API backend built in Go, serving a personal micro-frontend (MFE) shell architecture. Currently powers a job-search tracker (including recruiter management); planned expansions include vehicle maintenance, budgeting, project management, and recipes.

Built as both a learning project and a genuinely useful tool.

## Stack

- **Go 1.24**
- **[chi](https://github.com/go-chi/chi)** — lightweight HTTP router
- **MongoDB** — primary data store
- **JWT** — stateless access tokens with refresh token rotation
- **argon2id** — password hashing

## Architecture

Each domain follows a consistent three-layer structure:

```
<domain>/
  <domain>_handler.go   # HTTP request handling
  <domain>_routes.go    # Route registration (returns chi.Router)
  <domain>_store.go     # MongoDB data access
  <domain>_model.go     # Structs and request types
```

The `thing` package is the minimal reference implementation for this pattern. The `token` package adds a service layer for JWT business logic (minting, validation, refresh rotation).

### Auth Flow

- Access tokens are short-lived JWTs validated per-request via middleware
- Refresh token ids are stored in MongoDB with a TTL index for automatic expiry
- Refresh token rotation: each use revokes the old token and issues a new pair
- Passwords hashed with argon2id

## Running Locally

**Prerequisites:** Go 1.24+, a running MongoDB instance

**Environment:** Create a `.env` file in the repo root with the following keys:

```
MONGO_USERNAME=
MONGO_PASSWORD=
MONGO_CLUSTER=
MONGO_DB_NAME=
PRIVATE_KEY_PEM=
PUBLIC_KEY_PEM=
REFRESH_TOKEN_PEPPER=
```

**Run:**

```bash
go run main.go
```

Server listens on `localhost:8080`.

## Adding a New Service

Use `thing` as the structural template:

1. Create a `<domain>/` package with handler, routes, store, and model files
2. Add a store field to `app.Application` in `app/app.go`
3. Create the MongoDB collection and store instance in `main.go`
4. Register routes in `server/routes.go`

For auth-protected services that need the current user's ID, see `job` or `recruiter` as reference.

## Roadmap

- [x] Auth (JWT access + refresh token rotation)
- [x] Users
- [x] Job search tracker (with recruiter management)
- [ ] Vehicle maintenance tracker
- [ ] Budget tracker
- [ ] Project manager
- [ ] Recipes
