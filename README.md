<div align="center">

# homelab-status

**FastAPI service that exposes the status of services monitored by Uptime Kuma**

[![Deploy](https://github.com/ncorrea-13/homelab-status/actions/workflows/deploy.yml/badge.svg)](https://github.com/ncorrea-13/homelab-status/actions/workflows/deploy.yml)
[![Python](https://img.shields.io/badge/Python-3.12-3776AB?logo=python&logoColor=white)](https://www.python.org)
[![FastAPI](https://img.shields.io/badge/FastAPI-0.115-009688?logo=fastapi&logoColor=white)](https://fastapi.tiangolo.com)
[![SQLite](https://img.shields.io/badge/SQLite-file--based-003B57?logo=sqlite&logoColor=white)](https://www.sqlite.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

*[Versión en español](README.es.md)*

</div>

---

Receives notifications from [Uptime Kuma](https://github.com/louislam/uptime-kuma) (Webhook-type notification), stores the history in SQLite, and exposes the current status of each service as JSON. It's currently consumed by the homelab dashboard at [homelab.ncorrea.com.ar](https://homelab.ncorrea.com.ar).

Meant to run as one more container in the homelab (Podman), without exposing anything new to the internet beyond what you already choose to expose via Tunnel.

## Stack

| Layer | Technology |
| --- | --- |
| Runtime | Python 3.12 |
| Framework | FastAPI |
| ASGI server | Uvicorn |
| Database | SQLite (local file, mounted as a volume) |
| Container | Podman / Docker Compose |

## Quick Start

### Local (no container)

```bash
# 1. Install dependencies
pip install -r requirements.txt

# 2. Configure environment
cp .env.example .env
# Fill in WEBHOOK_SECRET

# 3. Run
uvicorn app.main:app --host 0.0.0.0 --port 8000
```

### Podman / Docker Compose

```bash
cp .env.example .env
# Fill in WEBHOOK_SECRET

podman compose up -d --build
# or: docker compose up -d --build
```

SQLite data persists in `status-data` (`/data/status.db` inside the container).

## Environment variables

| Variable | Required | Description |
| --- | --- | --- |
| `WEBHOOK_SECRET` | Yes | Sent as `x-webhook-secret` header or `?secret=` query param. Protects the Kuma webhook and the `/admin/services` endpoints. Empty = always 401. |
| `ALLOWED_ORIGINS` | No | Allowed CORS hosts (`app/main.py`). |
| `DB_PATH` | No | Path to the SQLite file (default `/data/status.db`). |

## API Endpoints

| Method | Route | Auth | Description |
| --- | --- | --- | --- |
| `POST` | `/webhook/kuma` | secret | Receives Uptime Kuma heartbeats and stores them in `events` |
| `GET` | `/api/status` | public | Returns the latest status of each active service |
| `POST` | `/admin/services` | secret | Create/edit a service (`{ id, name }`) |
| `DELETE` | `/admin/services/{id}` | secret | Soft-delete a service (`active = 0`) |
| `GET` | `/admin/services` | secret | Raw list of all services (debug) |

## Project Structure

```
app/
└── main.py           # FastAPI app: routes, CORS, SQLite access
compose.yml            # Service definition + status-data volume
Dockerfile              # Container image
requirements.txt        # Python dependencies
```

## Uptime Kuma integration

In any monitor, add a Webhook-type notification with the following URL:

```
https://<your-host>/webhook/kuma?secret=<WEBHOOK_SECRET>
```

(or send the secret as a header instead of a query param).

## Project

Personal service for Nicolás Correa's homelab.
