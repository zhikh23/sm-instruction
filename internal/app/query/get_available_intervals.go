package query

import (
	"context"
	"log/slog"
	"time"

	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type GetAvailableIntervals struct {
	ChatID       int64
	LocationUUID string
}

type GetAvailableIntervalsHandler decorator.QueryHandler[GetAvailableIntervals, []time.Time]

type getAvailableIntervalsHandler struct {
	chars sm.CharactersRepository
	locs  sm.LocationsRepository
}

func NewGetAvailableIntervalsHandler(
	chars sm.CharactersRepository,
	locs sm.LocationsRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetAvailableIntervalsHandler {
	if chars == nil {
		panic("characters repository is nil")
	}

	if locs == nil {
		panic("locations repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetAvailableIntervals, []time.Time](
		&getAvailableIntervalsHandler{chars: chars, locs: locs},
		log,
		metricsClient,
	)
}

func (h *getAvailableIntervalsHandler) Handle(ctx context.Context, query GetAvailableIntervals) ([]time.Time, error) {
	char, err := h.chars.Character(ctx, query.ChatID)
	if err != nil {
		return nil, err
	}

	finishAt, err := char.FinishTime()
	if err != nil {
		return nil, err
	}

	loc, err := h.locs.Location(ctx, query.LocationUUID)
	if err != nil {
		return nil, err
	}

	return loc.AvailableTimes(finishAt), nil
}
