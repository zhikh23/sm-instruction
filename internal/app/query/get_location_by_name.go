package query

import (
	"context"
	"log/slog"
	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type GetLocationByName struct {
	Name string
}

type GetLocationByNameHandler decorator.QueryHandler[GetLocationByName, Location]

type getLocationByNameHandler struct {
	locs sm.LocationsRepository
}

func NewGetLocationByNameHandler(
	locs sm.LocationsRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetLocationByNameHandler {
	if locs == nil {
		panic("locations repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetLocationByName, Location](
		&getLocationByNameHandler{locs: locs},
		log,
		metricsClient,
	)
}

func (h *getLocationByNameHandler) Handle(ctx context.Context, query GetLocationByName) (Location, error) {
	loc, err := h.locs.LocationByName(ctx, query.Name)
	if err != nil {
		return Location{}, err
	}

	return convertLocationToApp(loc), nil
}
