# LeetDrill

## Project Purpose

LeetDrill is a self-hosted LeetCode practice tracker. It imports the LeetCode
problem catalog, captures accepted submissions through a browser extension, and
builds a daily spaced-repetition queue from due, unsolved, and weak problems.

## Directory Structure

- `cmd/server/` - HTTP server entry point, routes, auth flows, extension API,
  sync worker startup.
- `cmd/ingest/` - LeetCode catalog ingester; can dry-run or upsert problems,
  patterns, and problem-pattern links.
- `internal/auth/` - web sessions, extension tokens, password auth, Google OAuth.
- `internal/leetcode/` - LeetCode GraphQL/verdict client code.
- `internal/mailer/` - SMTP mailer for multi-user email flows.
- `internal/models/` - shared domain types.
- `internal/srs/` - spaced-repetition scheduling logic.
- `internal/store/` - Postgres access layer using raw `pgx`; helpers accept
  `DBTX` so callers can use either the pool or a transaction.
- `internal/sync/` - cold-start history import and periodic recent-AC sync.
- `internal/vault/` - AES-GCM encryption for stored LeetCode cookies.
- `internal/web/` - embedded HTML templates, partials, favicon/static PNGs, and
  renderer helpers. The binary embeds these files with `go:embed`.
- `migrations/` - Goose SQL migrations.
- `extension/` - Chrome MV3 extension plus mirrored Firefox MV2 files under
  `extension/firefox/`.
- `scripts/` - extension validation/packaging and deployment scripts.
- `assets/`, `docs/`, `dist/` - repo assets, screenshots/docs, and generated
  build/package output.

## Build, Run, Test

Prereqs: Go 1.25+, Docker, and Task.

```sh
cp .env.example .env
# Fill LEETDRILL_COOKIE_KEY with: openssl rand -base64 32
task install:tools
task db:up
task migrate:up
task test
task dev
```

Common tasks:

- `task dev` - run the server with `air` if installed, otherwise `go run ./cmd/server`.
- `task build` - build `bin/server` and `bin/ingest`.
- `task test` - run `go test ./...`.
- `task vet` - run `go vet ./...`.
- `task migrate:up` / `task migrate:down` / `task migrate:status` - manage DB migrations.
- `task ingest -- -dry-run` - fetch LeetCode catalog pages without DB writes.
- `task ingest` - upsert the LeetCode catalog into the configured DB.
- `task extension:package` - validate and package Chrome/Firefox extension builds into `dist/extension-share/`.

The local Postgres service is defined in `docker-compose.yml` and listens on
host port `5433`. The default DSN is:

```text
postgres://leetdrill:leetdrill@localhost:5433/leetdrill?sslmode=disable
```

## Configuration

Runtime config comes from `.env` via `Taskfile.yml` or the process environment.
Required for the server:

- `DATABASE_URL`
- `LEETDRILL_COOKIE_KEY` - base64-encoded 32-byte key.

Important optional settings:

- `LEETDRILL_ADDR` - defaults to `:8080`.
- `LEETDRILL_BASE_PATH` - path prefix for reverse-proxy deployments.
- `LEETDRILL_APP_BASE` - public base URL for links.
- `SINGLE_USER=true` - self-host mode; ensures one user and can skip mailer setup.
- `USER_EMAIL`, `LEETCODE_USERNAME` - single-user bootstrap fields.
- `LEETDRILL_SYNC_WORKER=false` - disables the 30-minute recent-AC sync worker.
- `LEETDRILL_SECURE_COOKIES=true` - mark auth cookies secure.
- `SMTP_*`, `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` - multi-user email and Google login.

## Conventions

- Keep Go code formatted with `gofmt`; use standard `go test ./...` before
  handing off changes.
- Prefer the existing raw `pgx` store style. Store helpers should take
  `context.Context`, a `store.DBTX`, and explicit args; use `Store.InTx` when a
  flow needs atomic writes.
- Add schema changes as Goose migrations in `migrations/`; do not edit applied
  migration history unless intentionally rewriting local-only history.
- HTML templates are server-rendered with `html/template`; expose data through
  `web.PageData.Data` and register template helpers in `internal/web/render.go`.
- Use `web.AppPath`/the `appPath` template helper for links so
  `LEETDRILL_BASE_PATH` deployments keep working.
- Keep extension Chrome and Firefox shared files mirrored. After editing shared
  extension JS/HTML, update the matching `extension/firefox/` files and run
  `task extension:package`; `scripts/check_extensions.py` enforces many of
  these invariants.

## Gotchas

- The server requires `LEETDRILL_COOKIE_KEY` even in single-user mode because
  LeetCode cookies are encrypted before storage.
- `cmd/ingest` requires `DATABASE_URL` unless `-dry-run` is passed.
- `task migrate:*` uses `go run github.com/pressly/goose/v3/cmd/goose@latest`,
  so the first run may download the Goose tool.
- The periodic sync worker starts by default; set `LEETDRILL_SYNC_WORKER=false`
  for deterministic local debugging.
- `internal/web` is the embedded UI source used by the server. The top-level
  `web/` directory exists but is currently empty; the `tailwind` task still
  references `web/static/app.src.css`, so verify that path before relying on it.
- Extension store packages must not request localhost permissions; local
  backend testing may need temporary manifest edits that should not ship.
