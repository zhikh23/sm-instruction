package command

import (
	"context"
	"log/slog"
	"time"

	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type TakeSlot struct {
	GroupName    string
	ActivityName string
	Start        time.Time
}

type TakeSlotHandler decorator.CommandHandler[TakeSlot]

type takeSlotHandler struct {
	chars      sm.CharactersRepository
	activities sm.ActivitiesRepository
}

func NewTakeSlotHandler(
	chars sm.CharactersRepository,
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) TakeSlotHandler {
	if chars == nil {
		panic("characters repository is nil")
	}

	if activities == nil {
		panic("activities repository is nil")
	}

	return decorator.ApplyCommandDecorators[TakeSlot](
		&takeSlotHandler{chars, activities},
		log, metricsClient,
	)
}

func (h *takeSlotHandler) Handle(ctx context.Context, cmd TakeSlot) error {
	return h.chars.Update(
		ctx,
		cmd.GroupName,
		func(innerCtx1 context.Context, char *sm.Character) error {
			return h.activities.Update(
				innerCtx1,
				cmd.ActivityName,
				func(innerCtx2 context.Context, activity *sm.Activity) error {
					err := char.TakeSlot(cmd.Start, cmd.ActivityName)
					if err != nil {
						return err
					}
					return activity.TakeSlot(cmd.Start, cmd.GroupName)
				})
		})
}
