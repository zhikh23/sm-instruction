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
	activities sm.ActivitiesRepository
}

func NewGetLocationByNameHandler(
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetLocationByNameHandler {
	if activities == nil {
		panic("activities repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetLocationByName, Location](
		&getLocationByNameHandler{activities},
		log,
		metricsClient,
	)
}

func (h *getLocationByNameHandler) Handle(ctx context.Context, query GetLocationByName) (Location, error) {
	loc, err := h.activities.ActivityByName(ctx, query.Name)
	if err != nil {
		return Location{}, err
	}

	return convertLocationToApp(loc), nil
}
