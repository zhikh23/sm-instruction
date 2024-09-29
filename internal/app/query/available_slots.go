package query

import (
	"context"
	"log/slog"
	"time"

	"github.com/zhikh23/sm-instruction/internal/common/decorator"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
	"github.com/zhikh23/sm-instruction/pkg/funcs"
)

type AvailableSlots struct {
	GroupName    string
	ActivityName string
}

type AvailableSlotsHandler decorator.QueryHandler[AvailableSlots, []Slot]

type availableSlotsHandler struct {
	chars      sm.CharactersRepository
	activities sm.ActivitiesRepository
}

func NewAvailableSlotsHandler(
	chars sm.CharactersRepository,
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) AvailableSlotsHandler {
	if activities == nil {
		panic("activities repository is nil")
	}

	return decorator.ApplyQueryDecorators[AvailableSlots, []Slot](
		&availableSlotsHandler{chars, activities},
		log, metricsClient,
	)
}

func (h *availableSlotsHandler) Handle(ctx context.Context, query AvailableSlots) ([]Slot, error) {
	activity, err := h.activities.Activity(ctx, query.ActivityName)
	if err != nil {
		return nil, err
	}
	activitySlots := activity.AvailableSlots()

	char, err := h.chars.Character(ctx, query.GroupName)
	if err != nil {
		return nil, err
	}
	charSlots := char.AvailableSlots()

	availableSlots := sm.SlotsIntersection(activitySlots, charSlots)
	availableSlots = funcs.Filter(availableSlots, func(slot *sm.Slot) bool {
		return slot.Start.Before(*char.EndTime()) && slot.Start.After(time.Now().Add(-sm.MinDurationBefore))
	})

	return convertSlotsToApp(availableSlots), nil
}
