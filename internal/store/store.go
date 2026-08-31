package store

import (
	"context"

	"github.com/ncorrea-13/homelab-status/internal/models"
)

type Store interface {
	CreateEvent(ctx context.Context, eve models.Event) error
	CreateService(ctx context.Context, svc models.Service) error
	RemoveService(ctx context.Context, id int64) error
	GetServices(ctx context.Context) ([]models.Service, error)
	GetStatus(ctx context.Context) ([]models.ServiceStatus, error)
}
