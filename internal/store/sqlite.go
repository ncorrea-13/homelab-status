package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ncorrea-13/homelab-status/internal/models"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

var ErrServiceNotFound = errors.New("service not found")

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

func (s *SQLiteStore) CreateEvent(ctx context.Context, eve models.Event) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO events
			(monitor_id, monitor_name, monitor_url, status, message, duration_sec, checked_at, received_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		eve.MonitorID, eve.MonitorName, eve.MonitorURL, eve.Status, eve.Message, eve.DurationSec, eve.CheckedAt, eve.ReceivedAt,
	)
	return err
}

func (s *SQLiteStore) CreateService(ctx context.Context, svc models.Service) error {
	_, err := s.db.ExecContext(ctx, `
    INSERT INTO services (id, name, active) VALUES (?, ?, 1)
    ON CONFLICT(id) DO UPDATE SET name = excluded.name, active = 1
	`,
		svc.ID, svc.Name,
	)
	return err
}

func (s *SQLiteStore) RemoveService(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, "UPDATE services SET active = 0 WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrServiceNotFound
	}
	return nil
}

func (s *SQLiteStore) GetServices(ctx context.Context) ([]models.Service, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, active FROM services WHERE active = 1 ORDER BY name")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	services := []models.Service{}
	for rows.Next() {
		var svc models.Service
		if err := rows.Scan(&svc.ID, &svc.Name, &svc.Active); err != nil {
			return nil, err
		}
		services = append(services, svc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return services, nil
}

func (s *SQLiteStore) GetStatus(ctx context.Context) ([]models.ServiceStatus, error) {
	rows, err := s.db.QueryContext(ctx, `            
		SELECT s.id, s.name, e.status, e.message, e.checked_at
    FROM services s LEFT JOIN (
				SELECT e1.*
				FROM events e1
				INNER JOIN (
        		SELECT monitor_id, MAX(id) AS max_id 
						FROM events GROUP BY monitor_id
				) latest ON e1.monitor_id = latest.monitor_id AND e1.id = latest.max_id
		) e ON e.monitor_id = s.id
    WHERE s.active = 1
    ORDER BY s.name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	statuses := []models.ServiceStatus{}

	for rows.Next() {
		var srs models.ServiceStatus
		if err := rows.Scan(&srs.ID, &srs.Name, &srs.Status, &srs.Message, &srs.CheckedAt); err != nil {
			return nil, err
		}
		statuses = append(statuses, srs)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return statuses, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
