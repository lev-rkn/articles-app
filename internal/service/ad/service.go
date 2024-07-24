package ad

import (
	"ads-service/internal/store"
	"log/slog"

	def "ads-service/internal/service"
)

var _ def.AdService = (*service)(nil)

type service struct {
	store *store.Store
	logger *slog.Logger
}

func NewService(
	store *store.Store,
	logger *slog.Logger,
) *service {
	return &service{
		store: store,
		logger: logger,
	}
}
