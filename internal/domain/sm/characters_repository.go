package sm

import (
	"context"
	"errors"
)

var ErrCharacterAlreadyExists = errors.New("character already exists")
var ErrCharacterNotFound = errors.New("character not found")

type CharactersRepository interface {
	Save(ctx context.Context, character *Character) error
	Character(ctx context.Context, username string) (*Character, error)
	CharacterByGroup(ctx context.Context, groupName string) (*Character, error)
	Update(
		ctx context.Context,
		username string,
		updateFn func(innerCtx context.Context, char *Character) error,
	) error
}
