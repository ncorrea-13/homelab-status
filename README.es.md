<div align="center">

# homelab-status

**Servicio en Go que expone el estado de servicios monitoreados por Uptime Kuma**

[![Deploy](https://github.com/ncorrea-13/homelab-status/actions/workflows/ci.yml/badge.svg)](https://github.com/ncorrea-13/homelab-status/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![SQLite](https://img.shields.io/badge/SQLite-file--based-003B57?logo=sqlite&logoColor=white)](https://www.sqlite.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

_[English version](README.md)_

</div>

---

Recibe notificaciones webhook de [Uptime Kuma](https://github.com/louislam/uptime-kuma), guarda el historial en SQLite y expone el estado actual de cada servicio como JSON. Consumido por [homelab.ncorrea.com.ar](https://homelab.ncorrea.com.ar). Corre como un contenedor más del homelab (Podman), nada expuesto más allá de lo ya tuneleado.

## Quick Start

```bash
cp .env.example .env   # completar WEBHOOK_SECRET (o usar el secret de compose abajo)

go run ./cmd/server                 # local
podman compose up -d --build        # contenedor (necesita un secret podman webhook_secret)
```

Datos persisten en `/data/status.db` (SQLite, sin cgo).

## Variables de entorno

| Variable                                 | Descripción                                                                    |
| ------------------------------------------ | --------------------------------------------------------------------------------- |
| `WEBHOOK_SECRET` / `WEBHOOK_SECRET_FILE` | Una de las dos, requerida. Header `X-Webhook-Secret` o `?secret=`. Vacío = 401. |
| `ALLOWED_ORIGIN`                         | Origen CORS para `/api/status`.                                                |
| `DB_PATH`                                | Ruta SQLite (default `/data/status.db`).                                       |
| `PORT`                                   | Puerto host, solo compose.                                                     |

## API

| Método   | Ruta                    | Auth   |
| -------- | ----------------------- | ------ |
| `GET`    | `/healthz`              | -      |
| `GET`    | `/api/status`           | -      |
| `POST`   | `/webhook/kuma`         | secret |
| `GET`    | `/admin/services`       | secret |
| `POST`   | `/admin/services`       | secret |
| `DELETE` | `/admin/services/{id}`  | secret |

## Estructura

```
cmd/
├── server/       # entrypoint
└── healthcheck/  # binario del HEALTHCHECK
internal/
├── auth/         # carga del secret + middleware de auth
├── handlers/     # router, handlers HTTP, CORS
├── models/       # tipos de dominio
└── store/        # almacenamiento SQLite
```

## Integración con Uptime Kuma

Agregá una notificación Webhook apuntando a:

```
https://<tu-host>/webhook/kuma?secret=<WEBHOOK_SECRET>
```

## Proyecto

Servicio personal para el homelab de Nicolás Correa.
