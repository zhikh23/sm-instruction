package query

import (
	"context"
	"log/slog"
	"math"
	"slices"

	"github.com/zhikh23/sm-instruction/internal/common/decorator"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

type Rating struct {
}

type RatingHandler decorator.QueryHandler[Rating, []Character]

type ratingHandler struct {
	chars sm.CharactersRepository
}

func NewRatingHandler(
	chars sm.CharactersRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) RatingHandler {
	if chars == nil {
		panic("chars repository is nil")
	}

	return decorator.ApplyQueryDecorators[Rating, []Character](
		&ratingHandler{chars},
		log, metricsClient,
	)
}

func (h *ratingHandler) Handle(ctx context.Context, _ Rating) ([]Character, error) {
	chars, err := h.chars.Characters(ctx)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(chars, func(a, b *sm.Character) int {
		diff := b.Rating() - a.Rating()
		if diff == 0.0 {
			return 0
		}
		return int(math.Ceil(diff / math.Abs(diff))) // Получение знака
	})

	return convertCharactersToApp(chars), nil
}
