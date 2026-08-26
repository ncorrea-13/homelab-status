import os
import sqlite3
from contextlib import contextmanager
from datetime import datetime, timezone

from fastapi import FastAPI, Header, HTTPException, Query, Request
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel

DB_PATH = os.environ.get("DB_PATH", "/data/status.db")
WEBHOOK_SECRET = os.environ.get("WEBHOOK_SECRET", "")

ALLOWED_ORIGINS = [
    "https://homelab.ncorrea.com.ar",
]

STATUS_LABELS = {0: "down", 1: "up", 2: "pending", 3: "maintenance"}

app = FastAPI(title="homelab-status")

app.add_middleware(
    CORSMiddleware,
    allow_origins=ALLOWED_ORIGINS,
    allow_methods=["GET", "OPTIONS"],
    allow_headers=["*"],
)


@contextmanager
def get_db():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    try:
        yield conn
        conn.commit()
    finally:
        conn.close()


def init_db():
    with get_db() as conn:
        conn.execute(
            """
            CREATE TABLE IF NOT EXISTS services (
              id INTEGER PRIMARY KEY,
              name TEXT NOT NULL,
              active INTEGER NOT NULL DEFAULT 1
            )
            """
        )
        conn.execute(
            """
            CREATE TABLE IF NOT EXISTS events (
              id INTEGER PRIMARY KEY AUTOINCREMENT,
              monitor_id INTEGER NOT NULL,
              monitor_name TEXT NOT NULL,
              monitor_url TEXT,
              status TEXT NOT NULL,
              message TEXT,
              duration_sec INTEGER,
              checked_at TEXT NOT NULL,
              received_at TEXT NOT NULL
            )
            """
        )
        conn.execute(
            "CREATE INDEX IF NOT EXISTS idx_events_monitor_id ON events (monitor_id, id DESC)"
        )


@app.on_event("startup")
def on_startup():
    os.makedirs(os.path.dirname(DB_PATH), exist_ok=True)
    init_db()


def check_secret(provided: str | None):
    if not WEBHOOK_SECRET or provided != WEBHOOK_SECRET:
        raise HTTPException(status_code=401, detail="Unauthorized")


@app.post("/webhook/kuma")
async def webhook_kuma(
    request: Request,
    x_webhook_secret: str | None = Header(default=None),
    secret: str | None = Query(default=None),
):
    check_secret(x_webhook_secret or secret)

    payload = await request.json()
    monitor = payload.get("monitor")
    heartbeat = payload.get("heartbeat")

    if not monitor or not heartbeat:
        return {"status": "ok", "note": "test ping, nothing to store"}

    with get_db() as conn:
        conn.execute(
            """
            INSERT INTO events
              (monitor_id, monitor_name, monitor_url, status, message, duration_sec, checked_at, received_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
            """,
            (
                monitor["id"],
                monitor["name"],
                monitor.get("url"),
                STATUS_LABELS.get(heartbeat.get("status"), "unknown"),
                heartbeat.get("msg", ""),
                heartbeat.get("duration"),
                heartbeat.get("time") or datetime.now(timezone.utc).isoformat(),
                datetime.now(timezone.utc).isoformat(),
            ),
        )

    return {"status": "ok"}


class ServiceIn(BaseModel):
    id: int
    name: str


@app.post("/admin/services")
def add_service(
    service: ServiceIn,
    x_webhook_secret: str | None = Header(default=None),
    secret: str | None = Query(default=None),
):
    check_secret(x_webhook_secret or secret)

    with get_db() as conn:
        conn.execute(
            """
            INSERT INTO services (id, name, active) VALUES (?, ?, 1)
            ON CONFLICT(id) DO UPDATE SET name = excluded.name, active = 1
            """,
            (service.id, service.name),
        )

    return {"status": "ok"}


@app.delete("/admin/services/{service_id}")
def remove_service(
    service_id: int,
    x_webhook_secret: str | None = Header(default=None),
    secret: str | None = Query(default=None),
):
    check_secret(x_webhook_secret or secret)

    with get_db() as conn:
        conn.execute("UPDATE services SET active = 0 WHERE id = ?", (service_id,))

    return {"status": "ok"}


@app.get("/admin/services")
def list_services_raw(
    x_webhook_secret: str | None = Header(default=None),
    secret: str | None = Query(default=None),
):
    check_secret(x_webhook_secret or secret)

    with get_db() as conn:
        rows = conn.execute("SELECT * FROM services ORDER BY name").fetchall()

    return [dict(r) for r in rows]


# --- API pública -------------------------------------------------------


@app.get("/api/status")
def get_status():
    with get_db() as conn:
        rows = conn.execute(
            """
            SELECT
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
            ORDER BY s.name ASC
            """
        ).fetchall()

    return [dict(r) for r in rows]
