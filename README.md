# homelab-status

Cloudflare Worker que recibe notificaciones de [Uptime Kuma](https://github.com/louislam/uptime-kuma) (webhook tipo "Webhook") y expone el estado actual de los servicios monitoreados como API JSON. No sirve HTML: el consumo visual vive en una landing aparte que hace `fetch` a `/api/status`.

## Stack

- Cloudflare Workers (JS, sin build step)
- D1 (SQLite en el edge) para persistir servicios y eventos

## Setup

1. Instalar dependencias:
   ```
   npm install
   ```
2. Copiar `wrangler.toml.example` a `wrangler.toml` y completar `database_id` con el de tu base D1:
   ```
   cp wrangler.toml.example wrangler.toml
   npx wrangler d1 create homelab-status-db
   ```
3. Aplicar el schema:
   ```
   npx wrangler d1 execute homelab-status-db --file=schema.sql
   ```
4. Configurar el secret que valida el webhook y los endpoints de administración:
   ```
   npx wrangler secret put WEBHOOK_SECRET
   ```
5. Deploy:
   ```
   npx wrangler deploy
   ```

## Configuración

- **CORS**: el Worker solo responde con `Access-Control-Allow-Origin` para el host configurado en `ALLOWED_ORIGIN` (`src/index.js`). Cambiá ese valor por el dominio de tu landing.
- **WEBHOOK_SECRET**: se manda como header `x-webhook-secret` o query param `?secret=`. Protege el webhook de Kuma y los endpoints `/admin/services`.

## Endpoints

| Método | Ruta | Auth | Descripción |
|---|---|---|---|
| POST | `/webhook/kuma` | secret | Recibe heartbeats de Uptime Kuma y los guarda en `events` |
| GET | `/api/status` | pública | Devuelve el último estado de cada servicio activo |
| POST | `/admin/services` | secret | Alta/edición de un servicio (`{ id, name }`) |
| DELETE | `/admin/services?id=<id>` | secret | Baja lógica de un servicio (`active = 0`) |
| GET | `/admin/services` | secret | Lista cruda de todos los servicios |

## Uptime Kuma

En cada monitor: Notifications → Setup Notification → tipo **Webhook**, URL `https://<tu-worker>/webhook/kuma?secret=<WEBHOOK_SECRET>` (o mandar el secret por header).
