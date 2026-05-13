# leetdrill

Spaced-repetition layer on top of LeetCode. Web app + Chrome extension.

## Stack

Go 1.22 + chi + pgx + Postgres + htmx + Alpine + Tailwind. Manifest V3 extension.

## Layout

```
cmd/server         HTTP entrypoint
cmd/ingest         LeetCode catalog ingest CLI
internal/srs       SM-2 adapted scheduler (pure)
internal/leetcode  GraphQL client (public + authed)
internal/http      handlers, htmx fragments
internal/store     DB layer (sqlc target)
internal/auth      sessions, cookie vault wiring
internal/vault     AES-GCM for LeetCode cookies
internal/sync      periodic sync worker
internal/models    domain types
migrations/        goose SQL migrations
web/templates      html/template files
web/partials       htmx fragments
web/static         tailwind output, alpine, htmx
extension/         MV3 Chrome extension
```

## Bootstrap

Prereqs: Go 1.22+, Docker, [Task](https://taskfile.dev). Optional: sqlc, air, tailwindcss CLI.

```sh
cp .env.example .env
# fill LEETDRILL_COOKIE_KEY: openssl rand -base64 32

task install:tools     # goose, sqlc, air
task db:up             # start postgres on :5433
task migrate:up        # apply migrations
task test              # go test ./...
task dev               # serve on :8080
```

Smoke check:

```sh
curl localhost:8080/healthz
```

Then open:

- Web app: http://localhost:8080
- Today page: http://localhost:8080/session/today
- Extension options: `chrome://extensions` → Load unpacked → `extension/`

## Deployment

`scripts/deploy_server.sh` bootstraps or updates an Ubuntu-style VPS: packages,
source sync, `.env`, tests, build, migrations, systemd service, nginx proxy,
optional TLS via certbot, and extension share files.

Current production defaults:

```sh
task deploy:server -- abhiy.xyz
```

Useful first-run overrides:

```sh
LETSENCRYPT_EMAIL=you@example.com task deploy:server -- new-host.example.com

SETUP_POSTGRES=true \
DB_NAME=leetdrill \
DB_USER=leetdrill \
DB_PASSWORD='change-me' \
LETSENCRYPT_EMAIL=you@example.com \
task deploy:server -- new-host.example.com
```

By default the script uploads local `.env` only when the remote `.env` is
missing. Use `UPLOAD_ENV=true` to replace it, or `UPLOAD_ENV=false` to always
preserve or generate the remote `.env`.

Cold-start import:

- Web: Settings → Import LeetCode history.
- Extension: sync cookies, then popup → import history.

## Existing assets

`internal/srs/srs.go` + `srs_test.go` — SM-2 adapted scheduler, 4 ratings, leech detection, 24 subtests.
`internal/leetcode/client.go` — 6 GraphQL queries (4 public, 2 authed).

## Phases

0. Scaffold.
1. DB schema + ingest wiring.
2. Auth + cookie vault.
3. Extension MV3 + capture flow.
4. Web UI core (htmx pages).
5. Cold-start backfill + sync worker (current).
6. Retention: leech view, vacation, triage.
7. Stretch: PWA, multi-platform, judge, stats.
