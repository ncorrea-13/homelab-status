<div align="center">

# homelab-status

**API en FastAPI que expone el estado de servicios monitoreados por Uptime Kuma**

[![Deploy](https://github.com/ncorrea-13/homelab-status/actions/workflows/ci.yml/badge.svg)](https://github.com/ncorrea-13/homelab-status/actions/workflows/ci.yml)
[![Python](https://img.shields.io/badge/Python-3.12-3776AB?logo=python&logoColor=white)](https://www.python.org)
[![FastAPI](https://img.shields.io/badge/FastAPI-0.115-009688?logo=fastapi&logoColor=white)](https://fastapi.tiangolo.com)
[![SQLite](https://img.shields.io/badge/SQLite-file--based-003B57?logo=sqlite&logoColor=white)](https://www.sqlite.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

*[English version](README.md)*

</div>

---

Recibe notificaciones de [Uptime Kuma](https://github.com/louislam/uptime-kuma) (webhook tipo "Webhook"), guarda el historial en SQLite y expone el estado actual de cada servicio como JSON. Actualmente es consumida por la página del servidor en [homelab.ncorrea.com.ar](https://homelab.ncorrea.com.ar).

Pensado para correr como contenedor del homelab (Podman), sin exponer nada nuevo a internet salvo lo que ya decidas exponer vía Tunnel.

## Stack

| Capa | Tecnología |
| --- | --- |
| Runtime | Python 3.12 |
| Framework | FastAPI |
| Servidor ASGI | Uvicorn |
| Base de datos | SQLite (fichero local, montado como volumen) |
| Contenedor | Podman / Docker Compose |

## Quick Start

### Local (sin contenedor)

```bash
# 1. Instalar dependencias
pip install -r requirements.txt

# 2. Configurar entorno
cp .env.example .env
# Completar WEBHOOK_SECRET

# 3. Arrancar
uvicorn app.main:app --host 0.0.0.0 --port 8000
```

### Podman / Docker Compose

```bash
cp .env.example .env
# Completar WEBHOOK_SECRET

podman compose up -d --build
# o: docker compose up -d --build
```

SQLite persiste `status-data` (`/data/status.db` dentro del contenedor).

## Variables de entorno

| Variable | Requerida | Descripción |
| --- | --- | --- |
| `WEBHOOK_SECRET` | Sí | Header `x-webhook-secret` o query `?secret=`. Protege el webhook de Kuma y los endpoints `/admin/services`. Vacío = 401 siempre. |
| `ALLOWED_ORIGINS` | No | Hosts permitidos para CORS (`app/main.py`). |
| `DB_PATH` | No | Ruta del fichero SQLite (default `/data/status.db`). |

## API Endpoints

| Método | Ruta | Auth | Descripción |
| --- | --- | --- | --- |
| `POST` | `/webhook/kuma` | secret | Recibe heartbeats de Uptime Kuma y los guarda en `events` |
| `GET` | `/api/status` | pública | Devuelve el último estado de cada servicio activo |
| `POST` | `/admin/services` | secret | Alta/edición de un servicio (`{ id, name }`) |
| `DELETE` | `/admin/services/{id}` | secret | Baja lógica de un servicio (`active = 0`) |
| `GET` | `/admin/services` | secret | Lista cruda de todos los servicios (debug) |

## Estructura del Proyecto

```
app/
└── main.py           # App FastAPI: rutas, CORS, acceso a SQLite
compose.yml            # Definición del servicio + volumen status-data
Dockerfile              # Imagen del contenedor
requirements.txt        # Dependencias Python
```

## Integración con Uptime Kuma

En cualquier programa de monitoreo se debe agregar una notificación de tipo  Webhook, con la siguiente URL:

```
https://<tu-host>/webhook/kuma?secret=<WEBHOOK_SECRET>
```

(o mandar el secret por header en vez de query param).

## Proyecto

Servicio personal para el homelab de Nicolás Correa.
