package query

import (
	"context"
	"log/slog"
	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type GetCharacterByGroup struct {
	GroupName string
}

type GetCharacterByGroupHandler decorator.QueryHandler[GetCharacterByGroup, Character]

type getCharacterByGroupHandler struct {
	chars sm.CharactersRepository
}

func NewGetCharacterByGroupHandler(
	chars sm.CharactersRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetCharacterByGroupHandler {
	if chars == nil {
		panic("characters repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetCharacterByGroup, Character](
		&getCharacterByGroupHandler{chars: chars},
		log,
		metricsClient,
	)
}

func (h getCharacterByGroupHandler) Handle(ctx context.Context, query GetCharacterByGroup) (Character, error) {
	char, err := h.chars.CharacterByGroup(ctx, query.GroupName)
	if err != nil {
		return Character{}, err
	}

	return convertCharacterToApp(char), nil
}
