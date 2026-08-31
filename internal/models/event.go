package models

import "database/sql"

type Event struct {
	ID          int64          `json:"id"`
	MonitorID   int64          `json:"monitor_id"`
	MonitorName string         `json:"monitor_name"`
	MonitorURL  sql.NullString `json:"monitor_url"`
	Status      string         `json:"status"`
	Message     sql.NullString `json:"message"`
	DurationSec sql.NullInt64  `json:"duration_sec"`
	CheckedAt   string         `json:"checked_at"`
	ReceivedAt  string         `json:"received_at"`
}
