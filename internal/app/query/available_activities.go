package query

import (
	"context"
	"log/slog"

	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type AvailableActivities struct{}

type AvailableActivitiesHandler decorator.QueryHandler[AvailableActivities, []ActivityWithLocation]

type availableActivitiesHandler struct {
	activities sm.ActivitiesRepository
}

func NewAvailableActivitiesHandler(
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) AvailableActivitiesHandler {
	if activities == nil {
		panic("activities is nil")
	}

	return decorator.ApplyQueryDecorators[AvailableActivities, []ActivityWithLocation](
		&availableActivitiesHandler{activities},
		log, metricsClient,
	)
}

func (h *availableActivitiesHandler) Handle(
	ctx context.Context, _ AvailableActivities,
) ([]ActivityWithLocation, error) {
	activities, err := h.activities.ActivitiesWithLocations(ctx)
	if err != nil {
		return nil, err
	}
	return convertActivitiesWithLocations(activities), nil
}
