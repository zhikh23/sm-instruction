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
	activities sm.ActivitiesRepository
}

func NewGetLocationHandler(
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetLocationHandler {
	if activities == nil {
		panic("locations repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetLocation, Location](
		&getLocationHandler{activities},
		log,
		metricsClient,
	)
}

func (h *getLocationHandler) Handle(ctx context.Context, query GetLocation) (Location, error) {
	act, err := h.activities.Activity(ctx, query.UUID)
	if err != nil {
		return Location{}, err
	}

	_, err = act.LocationOrErr()
	if err != nil {
		return Location{}, err
	}

	return convertLocationToApp(act), nil
}
