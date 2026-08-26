# homelab-status

API en FastAPI que recibe notificaciones de [Uptime Kuma](https://github.com/louislam/uptime-kuma) (webhook tipo "Webhook"), guarda el historial en SQLite, y expone el estado actual de los servicios monitoreados como JSON. No sirve HTML: el consumo visual vive en una landing aparte que hace `fetch` a `/api/status`.

Pensado para correr como un contenedor más del homelab (Podman), sin exponer nada nuevo a internet salvo lo que ya decidas exponer via Tunnel.

## Stack

- FastAPI + Uvicorn
- SQLite (fichero local, montado como volumen en el contenedor)

## Setup

### Local (sin contenedor)

1. Instalar dependencias:
   ```
   pip install -r requirements.txt
   ```
2. Copiar `.env.example` a `.env` y completar `WEBHOOK_SECRET`:
   ```
   cp .env.example .env
   ```
3. Correr:
   ```
   uvicorn app.main:app --host 0.0.0.0 --port 8000
   ```

### Podman / Docker Compose

1. Copiar `.env.example` a `.env` y completar `WEBHOOK_SECRET`.
2. Levantar:
   ```
   podman compose up -d --build
   ```
   (o `docker compose up -d --build` si usás Docker)

La base SQLite persiste en el volumen `status-data` (`/data/status.db` dentro del contenedor).

## Configuración

- **CORS**: solo responde con `Access-Control-Allow-Origin` para los hosts en `ALLOWED_ORIGINS` (`app/main.py`). Cambiá esa lista por el dominio de tu landing.
- **WEBHOOK_SECRET**: se manda como header `x-webhook-secret` o query param `?secret=`. Protege el webhook de Kuma y los endpoints `/admin/services`. Si está vacío, esos endpoints quedan bloqueados (401 siempre).
- **DB_PATH**: ruta del fichero SQLite (default `/data/status.db`).

## Endpoints

| Método | Ruta | Auth | Descripción |
|---|---|---|---|
| POST | `/webhook/kuma` | secret | Recibe heartbeats de Uptime Kuma y los guarda en `events` |
| GET | `/api/status` | pública | Devuelve el último estado de cada servicio activo |
| POST | `/admin/services` | secret | Alta/edición de un servicio (`{ id, name }`) |
| DELETE | `/admin/services/{id}` | secret | Baja lógica de un servicio (`active = 0`) |
| GET | `/admin/services` | secret | Lista cruda de todos los servicios (debug) |

## Uptime Kuma

En cada monitor: Notifications → Setup Notification → tipo **Webhook**, URL `https://<tu-host>/webhook/kuma?secret=<WEBHOOK_SECRET>` (o mandar el secret por header).
