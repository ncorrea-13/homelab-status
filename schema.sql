-- Tabla de servicios: la lista fija de qué mostrar en el status page,
-- independiente de si Kuma ya mandó o no un evento para ellos.
CREATE TABLE IF NOT EXISTS services (
  id INTEGER PRIMARY KEY,           -- mismo id que monitor.id en Kuma
  name TEXT NOT NULL,
  active INTEGER NOT NULL DEFAULT 1 -- 1 = se muestra en el status page, 0 = oculto
);

-- Historial de eventos (heartbeats) reportados por Kuma vía webhook.
CREATE TABLE IF NOT EXISTS events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  monitor_id INTEGER NOT NULL,
  monitor_name TEXT NOT NULL,
  monitor_url TEXT,
  status TEXT NOT NULL,          -- 'up' | 'down' | 'pending' | 'maintenance' | 'unknown'
  message TEXT,
  duration_sec INTEGER,
  checked_at TEXT NOT NULL,
  received_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_monitor_id ON events (monitor_id, id DESC);
