package command

import (
	"context"
	"log/slog"

	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type RemoveBooking struct {
	Username string
}

type RemoveBookingHandler decorator.CommandHandler[RemoveBooking]

type removeBookingHandler struct {
	chars      sm.CharactersRepository
	activities sm.ActivitiesRepository
}

func NewRemoveBookingHandler(
	chars sm.CharactersRepository,
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) RemoveBookingHandler {
	if chars == nil {
		panic("characters repository is nil")
	}

	if activities == nil {
		panic("activities repository is nil")
	}

	return decorator.ApplyCommandDecorators[RemoveBooking](
		&removeBookingHandler{chars, activities},
		log,
		metricsClient,
	)
}

func (h *removeBookingHandler) Handle(ctx context.Context, cmd RemoveBooking) error {
	return h.chars.Update(ctx, cmd.Username, func(innerCtx context.Context, char *sm.Character) error {
		bookedTime, err := char.BookedTimeOrErr()
		if err != nil {
			return err
		}

		if err = char.RemoveBooking(); err != nil {
			return err
		}

		return h.activities.Update(
			innerCtx,
			bookedTime.ActivityUUID,
			func(innerCtx2 context.Context, activity *sm.Activity,
			) error {
				loc, err := activity.LocationOrErr()
				if err != nil {
					return err
				}

				return loc.RemoveBooking(bookedTime)
			})
	})
}
