package command

import (
	"context"
	"log/slog"
	"time"

	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type BookLocation struct {
	LocationUUID string
	ChatID       int64
	Time         time.Time
}

type BookLocationHandler decorator.CommandHandler[BookLocation]

type bookLocationHandler struct {
	chars sm.CharactersRepository
	locs  sm.LocationsRepository
}

func NewBookLocationHandler(
	chars sm.CharactersRepository,
	locs sm.LocationsRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) BookLocationHandler {
	if chars == nil {
		panic("characters repository is nil")
	}

	if locs == nil {
		panic("locations repository is nil")
	}

	return decorator.ApplyCommandDecorators[BookLocation](
		&bookLocationHandler{chars: chars, locs: locs},
		log,
		metricsClient,
	)
}

func (h *bookLocationHandler) Handle(ctx context.Context, cmd BookLocation) error {
	return h.locs.Update(ctx, cmd.LocationUUID, func(innerCtx context.Context, loc *sm.Location) error {
		return h.chars.Update(innerCtx, cmd.ChatID, func(innerCtx2 context.Context, char *sm.Character) error {
			return char.Book(loc, cmd.Time)
		})
	})
}
