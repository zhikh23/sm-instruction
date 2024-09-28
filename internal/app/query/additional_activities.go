package query

import (
	"context"
	"log/slog"

	"github.com/zhikh23/sm-instruction/internal/common/decorator"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

type AdditionalActivities struct {
	GroupName string
}

type AdditionalActivitiesHandler decorator.QueryHandler[AdditionalActivities, []Activity]

type additionalActivitiesHandler struct {
	activities sm.ActivitiesRepository
}

func NewAdditionalActivitiesHandler(
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) AdditionalActivitiesHandler {
	if activities == nil {
		panic("activities is nil")
	}

	return decorator.ApplyQueryDecorators[AdditionalActivities, []Activity](
		&additionalActivitiesHandler{activities},
		log, metricsClient,
	)
}

func (h *additionalActivitiesHandler) Handle(
	ctx context.Context, _ AdditionalActivities,
) ([]Activity, error) {
	activities, err := h.activities.AdditionalActivities(ctx)
	if err != nil {
		return nil, err
	}

	return convertActivitiesToApp(activities), nil
}
