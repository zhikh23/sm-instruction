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
	Username     string
	Time         time.Time
}

type BookLocationHandler decorator.CommandHandler[BookLocation]

type bookLocationHandler struct {
	chars      sm.CharactersRepository
	activities sm.ActivitiesRepository
}

func NewBookLocationHandler(
	chars sm.CharactersRepository,
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) BookLocationHandler {
	if chars == nil {
		panic("characters repository is nil")
	}

	if activities == nil {
		panic("activities repository is nil")
	}

	return decorator.ApplyCommandDecorators[BookLocation](
		&bookLocationHandler{chars, activities},
		log,
		metricsClient,
	)
}

func (h *bookLocationHandler) Handle(ctx context.Context, cmd BookLocation) error {
	bookTime, err := sm.NewBookedTime(cmd.Username, cmd.LocationUUID, cmd.Time, true)
	if err != nil {
		return err
	}

	return h.activities.Update(ctx, cmd.LocationUUID, func(innerCtx context.Context, activity *sm.Activity) error {
		loc, err := activity.LocationOrErr()
		if err != nil {
			return err
		}

		if err := loc.AddBooking(bookTime); err != nil {
			return err
		}

		return h.chars.Update(innerCtx, cmd.Username, func(innerCtx2 context.Context, char *sm.Character) error {
			return char.AddBooking(bookTime)
		})
	})
}
