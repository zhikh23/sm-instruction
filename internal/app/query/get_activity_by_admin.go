package query

import (
	"context"
	"log/slog"

	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type GetActivityByAdmin struct {
	Username string
}

type GetActivityByAdminHandler decorator.QueryHandler[GetActivityByAdmin, Activity]

type getActivityByAdmin struct {
	activities sm.ActivitiesRepository
}

func NewGetActivityByAdmin(
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetActivityByAdminHandler {
	if activities == nil {
		panic("activities repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetActivityByAdmin, Activity](
		&getActivityByAdmin{activities},
		log,
		metricsClient,
	)
}

func (h *getActivityByAdmin) Handle(ctx context.Context, query GetActivityByAdmin) (Activity, error) {
	act, err := h.activities.ActivityByAdmin(ctx, query.Username)
	if err != nil {
		return Activity{}, err
	}

	return convertActivityToApp(act), nil
}
