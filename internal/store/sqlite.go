package store

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Init(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return err
	}

	if _, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS services (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			active INTEGER NOT NULL DEFAULT 1
		)
	`); err != nil {
		return err
	}
	if _, err := s.db.ExecContext(ctx, `
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
	`); err != nil {
		return err
	}
	if _, err := s.db.ExecContext(ctx, `
     CREATE INDEX IF NOT EXISTS idx_events_monitor_id ON events (monitor_id, id DESC)
	`); err != nil {
		return err
	}
	return nil
}
