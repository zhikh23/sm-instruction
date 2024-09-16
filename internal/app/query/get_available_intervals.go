package query

import (
	"context"
	"log/slog"
	"time"

	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type GetAvailableIntervals struct {
	Username     string
	ActivityUUID string
}

type GetAvailableIntervalsHandler decorator.QueryHandler[GetAvailableIntervals, []time.Time]

type getAvailableIntervalsHandler struct {
	chars      sm.CharactersRepository
	activities sm.ActivitiesRepository
}

func NewGetAvailableIntervalsHandler(
	chars sm.CharactersRepository,
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetAvailableIntervalsHandler {
	if chars == nil {
		panic("characters repository is nil")
	}

	if activities == nil {
		panic("activities repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetAvailableIntervals, []time.Time](
		&getAvailableIntervalsHandler{chars, activities},
		log,
		metricsClient,
	)
}

func (h *getAvailableIntervalsHandler) Handle(ctx context.Context, query GetAvailableIntervals) ([]time.Time, error) {
	char, err := h.chars.Character(ctx, query.Username)
	if err != nil {
		return nil, err
	}

	finishAt, err := char.FinishTime()
	if err != nil {
		return nil, err
	}

	act, err := h.activities.Activity(ctx, query.ActivityUUID)
	if err != nil {
		return nil, err
	}

	loc, err := act.LocationOrErr()
	if err != nil {
		return nil, err
	}

	return loc.AvailableTimes(finishAt), nil
}
