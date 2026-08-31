package handlers

import (
	"github.com/ncorrea-13/homelab-status/internal/models"
)

func toNullString(s *string) models.NullString {
	if s == nil {
		return models.NullString{}
	}
	return models.NullString{String: *s, Valid: true}
}

func toNullInt64(i *int64) models.NullInt64 {
	if i == nil {
		return models.NullInt64{}
	}
	return models.NullInt64{Int64: *i, Valid: true}
}
