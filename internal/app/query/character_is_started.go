package query

import (
	"context"
	"errors"
	"log/slog"

	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type CharacterIsStarted struct {
	Username string
}

type CharacterIsStartedHandler decorator.QueryHandler[CharacterIsStarted, bool]

type characterIsStartedHandler struct {
	chars sm.CharactersRepository
}

func NewCharacterIsStarted(
	chars sm.CharactersRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) CharacterIsStartedHandler {
	if chars == nil {
		panic("characters repository is nil")
	}

	return decorator.ApplyQueryDecorators[CharacterIsStarted, bool](
		&characterIsStartedHandler{chars: chars},
		log,
		metricsClient,
	)
}

func (h *characterIsStartedHandler) Handle(ctx context.Context, query CharacterIsStarted) (bool, error) {
	char, err := h.chars.Character(ctx, query.Username)
	if errors.Is(err, sm.ErrCharacterNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return char.IsStarted(), nil
}
