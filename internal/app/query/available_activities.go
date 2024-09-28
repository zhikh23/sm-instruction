package query

import (
	"context"
	"github.com/zhikh23/sm-instruction/pkg/funcs"
	"log/slog"

	"github.com/zhikh23/sm-instruction/internal/common/decorator"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

type AvailableActivities struct {
	GroupName string
}

type AvailableActivitiesHandler decorator.QueryHandler[AvailableActivities, []Activity]

type availableActivitiesHandler struct {
	chars      sm.CharactersRepository
	activities sm.ActivitiesRepository
}

func NewAvailableActivitiesHandler(
	chars sm.CharactersRepository,
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) AvailableActivitiesHandler {
	if activities == nil {
		panic("activities is nil")
	}

	return decorator.ApplyQueryDecorators[AvailableActivities, []Activity](
		&availableActivitiesHandler{chars, activities},
		log, metricsClient,
	)
}

func (h *availableActivitiesHandler) Handle(
	ctx context.Context, q AvailableActivities,
) ([]Activity, error) {
	char, err := h.chars.Character(ctx, q.GroupName)
	if err != nil {
		return nil, err
	}

	activities, err := h.activities.AvailableActivities(ctx)
	if err != nil {
		return nil, err
	}

	activities = funcs.Filter(activities, func(act *sm.Activity) bool {
		return len(sm.SlotsIntersection(act.AvailableSlots(), char.AvailableSlots())) > 0 && !act.HasTaken(q.GroupName)
	})

	return convertActivitiesToApp(activities), nil
}
