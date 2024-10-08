package command

import (
	"context"
	"github.com/zhikh23/sm-instruction/internal/common/decorator"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
	"log/slog"
)

type AwardCharacter struct {
	GroupName    string
	ActivityName string
	SkillType    string
	Points       int
}

type AwardCharacterHandler decorator.CommandHandler[AwardCharacter]

type awardCharacterHandler struct {
	chars      sm.CharactersRepository
	activities sm.ActivitiesRepository
}

func NewAwardCharacterHandler(
	chars sm.CharactersRepository,
	activities sm.ActivitiesRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) AwardCharacterHandler {
	if chars == nil {
		panic("characters repository is nil")
	}

	if activities == nil {
		panic("activities repository is nil")
	}

	return decorator.ApplyCommandDecorators[AwardCharacter](
		&awardCharacterHandler{chars, activities},
		log, metricsClient,
	)
}

func (h *awardCharacterHandler) Handle(ctx context.Context, cmd AwardCharacter) error {
	st, err := sm.NewSkillTypeFromString(cmd.SkillType)
	if err != nil {
		return err
	}

	act, err := h.activities.Activity(ctx, cmd.ActivityName)
	if err != nil {
		return err
	}

	return h.chars.Update(ctx, cmd.GroupName, func(innerCtx context.Context, char *sm.Character) error {
		return act.Award(char, st, cmd.Points)
	})
}
