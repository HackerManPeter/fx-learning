# fx_tutorial

A compact tutorial application showing how to build a small web service using Go, GoFiber (FastHTTP), and Uber's Fx for dependency injection and lifecycle management.

This repository is intended for learning and experimenting with Fx wiring patterns, application structure, and simple service composition.

**Prerequisites**

- Go 1.20+ installed (GOMOD enabled)
- Docker (optional, for container builds)

**Quick start**

1.  From the repository root, download dependencies:

    ```bash
    go mod tidy
    ```

2.  Run the application locally:

    ```bash
    go run main.go
    ```

    The app bootstraps an `fx.App` and starts the HTTP server (check `internal/config` for the configured port; Fiber defaults are commonly `3000`).

3.  Build and run with Docker (optional):

    ```bash
    docker build -t fx-tutorial .
    docker run -p 3000:3000 fx-tutorial
    ```

**Useful commands**

- Download modules: `go mod tidy`
- Run locally: `go run main.go`
- Build: `go build -v ./...`

**Project structure**

- `main.go` — application entrypoint: constructs the Fx application, registers providers and invocations.
- `server/` — server bootstrap and wiring (`server.go`).
- `internal/config/` — configuration loader and environment wiring.
- `internal/database/` — database connection provider.
- `internal/cache/` — redis or in-memory cache provider.
- `internal/auth/` — authentication service, routes, and repository wiring.
- `internal/bible/` — example feature module (routes and service).
- `internal/middleware/` — middleware providers and wiring.
- `internal/models/` — domain models and response objects.

Files and folders follow a convention inspired by Go best practices and the tutorial's goal of keeping internals unexported.

**How the app is wired (overview)**

- Constructors (providers) are registered with `fx.Provide` in `main.go` (for example: `config.NewConfig`, `server.NewServer`, `database.NewDatabase`, `cache.NewCache`).
- Feature modules (auth, bible, middleware) expose service constructors and route registration helpers that are provided to the Fx graph.
- `fx.Invoke` is used to run start-up wiring code and ensure required components are constructed on start.

**Configuration & environment**

- Environment variables are loaded automatically via `godotenv` (see `main.go` import of `github.com/joho/godotenv/autoload`).
- See `internal/config` for the expected environment keys (port, database URL, redis URL, etc.).

**Notes & learnings**

- This project explores structuring an Fx application similar to patterns in frameworks like NestJS: small service constructors, explicit dependency injection, and route registration attached to service receivers.
- The `learnings.md` and `explanation.md` files include the author's reflections on FastHTTP vs `net/http`, the reasons for choosing receiver methods for handlers, and design trade-offs around repositories and testing.

**Contributing**

- Feel free to open issues or PRs for improvements, additional examples, or exercises.
- Keep changes small, and add tests for new behavior.

**Next steps / Suggestions**

- Add a `Makefile` or `scripts/` to automate common tasks (run, build, lint, test).
- Add tests
