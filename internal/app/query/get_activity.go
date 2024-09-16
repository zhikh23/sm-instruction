package query

import (
	"context"
	"log/slog"

	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type GetActivity struct {
	ActivityUUID string
}

type GetActivityHandler decorator.QueryHandler[GetActivity, Activity]

type getActivityHandler struct {
	activities sm.ActivitiesRepository
}

func NewGetActivityHandler(
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetActivityHandler {
	if activities == nil {
		panic("activities repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetActivity, Activity](
		&getActivityHandler{activities},
		log,
		metricsClient,
	)
}

func (h *getActivityHandler) Handle(ctx context.Context, query GetActivity) (Activity, error) {
	act, err := h.activities.Activity(ctx, query.ActivityUUID)
	if err != nil {
		return Activity{}, err
	}

	return convertActivityToApp(act), nil
}
