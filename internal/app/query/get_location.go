package query

import (
	"context"
	"log/slog"
	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type GetLocation struct {
	UUID string
}

type GetLocationHandler decorator.QueryHandler[GetLocation, Location]

type getLocationHandler struct {
	locs sm.LocationsRepository
}

func NewGetLocationHandler(
	locs sm.LocationsRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetLocationHandler {
	if locs == nil {
		panic("locations repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetLocation, Location](
		&getLocationHandler{locs: locs},
		log,
		metricsClient,
	)
}

func (h *getLocationHandler) Handle(ctx context.Context, query GetLocation) (Location, error) {
	loc, err := h.locs.Location(ctx, query.UUID)
	if err != nil {
		return Location{}, err
	}
	return convertLocationToApp(loc), nil
}
