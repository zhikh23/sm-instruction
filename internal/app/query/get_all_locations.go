package query

import (
	"context"
	"log/slog"

	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type GetAllLocations struct {
}

type GetAllLocationsHandler decorator.QueryHandler[GetAllLocations, []Location]

type getAllLocationsHandler struct {
	locs sm.LocationsRepository
}

func NewGetAllLocationsHandler(
	locs sm.LocationsRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetAllLocationsHandler {
	if locs == nil {
		panic("locations repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetAllLocations, []Location](
		&getAllLocationsHandler{locs: locs},
		log,
		metricsClient,
	)
}

func (h *getAllLocationsHandler) Handle(ctx context.Context, query GetAllLocations) ([]Location, error) {
	locs, err := h.locs.Locations(ctx)
	if err != nil {
		return nil, err
	}

	return convertLocationsToApp(locs), nil
}
