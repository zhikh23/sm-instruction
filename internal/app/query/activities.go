package query

import (
	"context"
	"github.com/zhikh23/sm-instruction/internal/common/decorator"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
	"log/slog"
)

type Activities struct {
}

type ActivitiesHandler decorator.QueryHandler[Activities, []Activity]

type activitiesHandler struct {
	activities sm.ActivitiesRepository
}

func NewActivitiesHandler(
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) ActivitiesHandler {
	if activities == nil {
		panic("activities repository is nil")
	}

	return decorator.ApplyQueryDecorators[Activities, []Activity](
		&activitiesHandler{activities},
		log, metricsClient,
	)
}

func (h *activitiesHandler) Handle(ctx context.Context, _ Activities) ([]Activity, error) {
	activities, err := h.activities.Activities(ctx)
	if err != nil {
		return nil, err
	}

	return convertActivitiesToApp(activities), nil
}
