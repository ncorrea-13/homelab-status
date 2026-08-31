package models

type ServiceStatus struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Status    NullString `json:"status"`
	Message   NullString `json:"message"`
	CheckedAt NullString `json:"checked_at"`
}
