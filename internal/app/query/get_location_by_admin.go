package query

import (
	"context"
	"log/slog"
	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type GetLocationByAdmin struct {
	Username string
}

type GetLocationByAdminHandler decorator.QueryHandler[GetLocationByAdmin, Location]

type getLocationByAdminHandler struct {
	locs sm.LocationsRepository
}

func NewGetLocationByAdminHandler(
	locs sm.LocationsRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetLocationByAdminHandler {
	if locs == nil {
		panic("locations repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetLocationByAdmin, Location](
		&getLocationByAdminHandler{locs: locs},
		log,
		metricsClient,
	)
}

func (h *getLocationByAdminHandler) Handle(ctx context.Context, query GetLocationByAdmin) (Location, error) {
	loc, err := h.locs.LocationByAdmin(ctx, query.Username)
	if err != nil {
		return Location{}, err
	}

	return convertLocationToApp(loc), nil
}
