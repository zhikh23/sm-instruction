package command

import (
	"context"
	"log/slog"
	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type CancelBooking struct {
	ChatID int64
}

type CancelBookingHandler decorator.CommandHandler[CancelBooking]

type cancelBookingHandler struct {
	chars sm.CharactersRepository
	locs  sm.LocationsRepository
}

func NewCancelBookingHandler(
	chars sm.CharactersRepository,
	locs sm.LocationsRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) CancelBookingHandler {
	if chars == nil {
		panic("characters repository is nil")
	}

	if locs == nil {
		panic("locations repository is nil")
	}

	return decorator.ApplyCommandDecorators[CancelBooking](
		&cancelBookingHandler{chars: chars, locs: locs},
		log,
		metricsClient,
	)
}

func (h *cancelBookingHandler) Handle(ctx context.Context, cmd CancelBooking) error {
	return h.chars.Update(ctx, cmd.ChatID, func(innerCtx context.Context, char *sm.Character) error {
		locUUID, err := char.BookedLocation()
		if err != nil {
			return err
		}

		return h.locs.Update(innerCtx, locUUID, func(innerCtx2 context.Context, loc *sm.Location) error {
			return char.RemoveBooking(loc)
		})
	})
}
