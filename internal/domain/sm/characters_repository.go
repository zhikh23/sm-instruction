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
	Update(
		ctx context.Context,
		chatID int64,
		updateFn func(innerCtx context.Context, char *Character) error,
	) error
}
