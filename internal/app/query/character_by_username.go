package query

import (
	"context"
	"log/slog"

	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type CharacterByUsername struct {
	Username string
}

type CharacterByUsernameHandler decorator.QueryHandler[CharacterByUsername, Character]

type getCharacterByUsernameHandler struct {
	chars sm.CharactersRepository
}

func NewCharacterByUsernameHandler(
	chars sm.CharactersRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) CharacterByUsernameHandler {
	if chars == nil {
		panic("characters repository is nil")
	}

	return decorator.ApplyQueryDecorators[CharacterByUsername, Character](
		&getCharacterByUsernameHandler{chars: chars},
		log,
		metricsClient,
	)
}

func (h getCharacterByUsernameHandler) Handle(ctx context.Context, query CharacterByUsername) (Character, error) {
	char, err := h.chars.CharacterByUsername(ctx, query.Username)
	if err != nil {
		return Character{}, err
	}

	return convertCharacterToApp(char), nil
}
