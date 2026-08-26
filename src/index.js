// Worker que recibe notificaciones de Uptime Kuma (tipo "Webhook") y
// expone el estado actual como API JSON. No sirve HTML - el consumo
// visual vive en la landing existente (homelab.html), que hace fetch
// a /api/status y renderiza con su propio CSS.
//
// Tablas en D1:
//   - "services": la lista fija de que mostrar. Se administra con
//     los endpoints /admin/services.
//   - "events": el historial de heartbeats que manda Kuma via webhook.

const STATUS_LABELS = {
  0: "down",
  1: "up",
  2: "pending",
  3: "maintenance",
};

// CORS: permití acceder desde tu dominio. Ajustá si tu landing vive
// en otro host.
const ALLOWED_ORIGIN = "https://homelab.ncorrea.com.ar";

function withCors(response) {
  const headers = new Headers(response.headers);
  headers.set("Access-Control-Allow-Origin", ALLOWED_ORIGIN);
  headers.set("Access-Control-Allow-Methods", "GET, OPTIONS");
  return new Response(response.body, { status: response.status, headers });
}

export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    const method = request.method;

    if (method === "OPTIONS") {
      return withCors(new Response(null, { status: 204 }));
    }

    // --- Webhook de Kuma ---
    if (url.pathname === "/webhook/kuma" && method === "POST") {
      return handleKumaWebhook(request, env);
    }

    // --- Administracion de servicios (protegida con el mismo secret) ---
    if (url.pathname === "/admin/services" && method === "POST") {
      return handleAddService(request, env);
    }
    if (url.pathname === "/admin/services" && method === "DELETE") {
      return handleRemoveService(request, env);
    }
    if (url.pathname === "/admin/services" && method === "GET") {
      return handleListServicesRaw(request, env);
    }

    // --- API publica (unica ruta consumida por la landing) ---
    if (url.pathname === "/api/status") {
      const current = await getCurrentStatuses(env);
      return withCors(Response.json(current));
    }

    return new Response("Not found", { status: 404 });
  },
};

function checkSecret(request, env) {
  const provided =
    request.headers.get("x-webhook-secret") ||
    new URL(request.url).searchParams.get("secret");
  return Boolean(env.WEBHOOK_SECRET) && provided === env.WEBHOOK_SECRET;
}

async function handleKumaWebhook(request, env) {
  if (!checkSecret(request, env)) {
    return new Response("Unauthorized", { status: 401 });
  }

  let payload;
  try {
    payload = await request.json();
  } catch {
    return new Response("Invalid JSON", { status: 400 });
  }

  const monitor = payload.monitor;
  const heartbeat = payload.heartbeat;

  if (!monitor || !heartbeat) {
    return new Response("OK (test ping, nothing to store)", { status: 200 });
  }

  await env.DB.prepare(
    `INSERT INTO events (monitor_id, monitor_name, monitor_url, status, message, duration_sec, checked_at, received_at)
     VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
  )
    .bind(
      monitor.id,
      monitor.name,
      monitor.url || null,
      STATUS_LABELS[heartbeat.status] ?? "unknown",
      heartbeat.msg || "",
      heartbeat.duration ?? null,
      heartbeat.time || new Date().toISOString(),
      new Date().toISOString(),
    )
    .run();

  return new Response("OK", { status: 200 });
}

async function handleAddService(request, env) {
  if (!checkSecret(request, env)) {
    return new Response("Unauthorized", { status: 401 });
  }

  let body;
  try {
    body = await request.json();
  } catch {
    return new Response("Invalid JSON", { status: 400 });
  }

  if (!body.id || !body.name) {
    return new Response("Missing id or name", { status: 400 });
  }

  await env.DB.prepare(
    `INSERT INTO services (id, name, active) VALUES (?, ?, 1)
     ON CONFLICT(id) DO UPDATE SET name = excluded.name, active = 1`,
  )
    .bind(body.id, body.name)
    .run();

  return new Response("OK", { status: 200 });
}

async function handleRemoveService(request, env) {
  if (!checkSecret(request, env)) {
    return new Response("Unauthorized", { status: 401 });
  }

  const id = new URL(request.url).searchParams.get("id");
  if (!id) {
    return new Response("Missing id", { status: 400 });
  }

  await env.DB.prepare(`UPDATE services SET active = 0 WHERE id = ?`)
    .bind(id)
    .run();

  return new Response("OK", { status: 200 });
}

async function handleListServicesRaw(request, env) {
  if (!checkSecret(request, env)) {
    return new Response("Unauthorized", { status: 401 });
  }
  const { results } = await env.DB.prepare(
    `SELECT * FROM services ORDER BY name`,
  ).all();
  return Response.json(results);
}

async function getCurrentStatuses(env) {
  const { results } = await env.DB.prepare(
    `SELECT
       s.id,
       s.name,
       e.status,
       e.message,
       e.checked_at
     FROM services s
     LEFT JOIN (
       SELECT e1.*
       FROM events e1
       INNER JOIN (
         SELECT monitor_id, MAX(id) AS max_id FROM events GROUP BY monitor_id
       ) latest ON e1.monitor_id = latest.monitor_id AND e1.id = latest.max_id
     ) e ON e.monitor_id = s.id
     WHERE s.active = 1
     ORDER BY s.name ASC`,
  ).all();

  return results;
}
