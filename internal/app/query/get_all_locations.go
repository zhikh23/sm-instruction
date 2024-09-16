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
	activities sm.ActivitiesRepository
}

func NewGetAllLocationsHandler(
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetAllLocationsHandler {
	if activities == nil {
		panic("activities repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetAllLocations, []Location](
		&getAllLocationsHandler{activities},
		log,
		metricsClient,
	)
}

func (h *getAllLocationsHandler) Handle(ctx context.Context, _ GetAllLocations) ([]Location, error) {
	acts, err := h.activities.ActivitiesWithLocations(ctx)
	if err != nil {
		return nil, err
	}

	return convertLocationsToApp(acts), nil
}
