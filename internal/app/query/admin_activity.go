package query

import (
	"context"
	"log/slog"

	"github.com/zhikh23/sm-instruction/internal/common/decorator"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

type AdminActivity struct {
	Username string
}

type AdminActivityHandler decorator.QueryHandler[AdminActivity, Activity]

type adminActivityHandler struct {
	activities sm.ActivitiesRepository
}

func NewAdminActivtyHandler(
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) AdminActivityHandler {
	if activities == nil {
		panic("activities repository is nil")
	}

	return decorator.ApplyQueryDecorators[AdminActivity, Activity](
		&adminActivityHandler{activities},
		log,
		metricsClient,
	)
}

func (h *adminActivityHandler) Handle(ctx context.Context, query AdminActivity) (Activity, error) {
	act, err := h.activities.ActivityByAdmin(ctx, query.Username)
	if err != nil {
		return Activity{}, err
	}

	return convertActivityToApp(act), nil
}
