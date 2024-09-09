package sm

import (
	"context"
	"errors"
)

var (
	ErrCharacterAlreadyExists = errors.New("character already exists")
	ErrCharacterNotFound      = errors.New("character not found")
)

type CharactersRepository interface {
	Save(ctx context.Context, character *Character) error
	Character(ctx context.Context, chatID int64) (*Character, error)
	SortedByRating(ctx context.Context) ([]*Character, error)
	Update(
		ctx context.Context,
		chatID int64,
		updateFn func(context.Context, *Character) error,
	) error
}
