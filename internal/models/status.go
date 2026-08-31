package models

import "database/sql"

type ServiceStatus struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Status    sql.NullString `json:"status"`
	Message   sql.NullString `json:"message"`
	CheckedAt sql.NullString `json:"checked_at"`
}
