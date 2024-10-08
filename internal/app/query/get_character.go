package query

import (
	"context"
	"log/slog"

	"github.com/zhikh23/sm-instruction/internal/common/decorator"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

type GetCharacter struct {
	GroupName string
}

type GetCharacterHandler decorator.QueryHandler[GetCharacter, Character]

type getCharacterHandler struct {
	chars sm.CharactersRepository
}

func NewGetCharacterHandler(
	chars sm.CharactersRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetCharacterHandler {
	if chars == nil {
		panic("characters repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetCharacter, Character](
		&getCharacterHandler{chars: chars},
		log,
		metricsClient,
	)
}

func (h getCharacterHandler) Handle(ctx context.Context, query GetCharacter) (Character, error) {
	char, err := h.chars.Character(ctx, query.GroupName)
	if err != nil {
		return Character{}, err
	}

	return convertCharacterToApp(char), nil
}
