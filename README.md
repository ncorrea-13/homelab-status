<div align="center">

# homelab-status

**Go service that exposes the status of services monitored by Uptime Kuma**

[![Deploy](https://github.com/ncorrea-13/homelab-status/actions/workflows/ci.yml/badge.svg)](https://github.com/ncorrea-13/homelab-status/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![SQLite](https://img.shields.io/badge/SQLite-file--based-003B57?logo=sqlite&logoColor=white)](https://www.sqlite.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

_[Versión en español](README.es.md)_

</div>

---

Receives [Uptime Kuma](https://github.com/louislam/uptime-kuma) webhook notifications, stores the history in SQLite, and exposes each service's current status as JSON. Consumed by [homelab.ncorrea.com.ar](https://homelab.ncorrea.com.ar). Runs as one more container in the homelab (Podman), nothing exposed beyond what's already tunneled.

## Quick Start

```bash
cp .env.example .env   # fill in WEBHOOK_SECRET (or use the compose secret below)

go run ./cmd/server                 # local
podman compose up -d --build        # container (needs a webhook_secret podman secret)
```

Data persists at `/data/status.db` (SQLite, no cgo).

## Environment variables

| Variable                            | Description                                                                  |
| ------------------------------------ | ------------------------------------------------------------------------------ |
| `WEBHOOK_SECRET` / `WEBHOOK_SECRET_FILE` | One of the two, required. Header `X-Webhook-Secret` or `?secret=`. Empty = 401. |
| `ALLOWED_ORIGIN`                     | CORS origin for `/api/status`.                                                |
| `DB_PATH`                            | SQLite path (default `/data/status.db`).                                      |
| `PORT`                               | Host port, compose only.                                                      |

## API

| Method   | Route                  | Auth   |
| -------- | ---------------------- | ------ |
| `GET`    | `/healthz`             | -      |
| `GET`    | `/api/status`          | -      |
| `POST`   | `/webhook/kuma`        | secret |
| `GET`    | `/admin/services`      | secret |
| `POST`   | `/admin/services`      | secret |
| `DELETE` | `/admin/services/{id}` | secret |

## Structure

```
cmd/
├── server/       # entrypoint
└── healthcheck/  # HEALTHCHECK binary
internal/
├── auth/         # secret loading + auth middleware
├── handlers/     # router, HTTP handlers, CORS
├── models/       # domain types
└── store/        # SQLite storage
```

## Uptime Kuma integration

Add a Webhook notification pointing to:

```
https://<your-host>/webhook/kuma?secret=<WEBHOOK_SECRET>
```

## Project

Personal service for Nicolás Correa's homelab.
