package mocks

import (
	"context"
	"sync"

	"sm-instruction/internal/domain/sm"
)

type mockCharactersRepository struct {
	m map[string]sm.Character
	sync.RWMutex
}

func NewMockCharactersRepository() sm.CharactersRepository {
	return &mockCharactersRepository{
		m: make(map[string]sm.Character),
	}
}

func (r *mockCharactersRepository) Save(_ context.Context, char *sm.Character) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[char.Username]; ok {
		return sm.ErrCharacterAlreadyExists
	}

	r.m[char.Username] = *char

	return nil
}

func (r *mockCharactersRepository) Character(_ context.Context, username string) (*sm.Character, error) {
	r.RLock()
	defer r.RUnlock()

	char, ok := r.m[username]
	if !ok {
		return nil, sm.ErrCharacterNotFound
	}

	return &char, nil
}

func (r *mockCharactersRepository) Update(
	ctx context.Context,
	username string,
	updateFn func(innerCtx context.Context, char *sm.Character) error,
) error {
	r.Lock()
	defer r.Unlock()

	char, ok := r.m[username]
	if !ok {
		return sm.ErrCharacterNotFound
	}

	err := updateFn(ctx, &char)
	if err != nil {
		return err
	}

	r.m[username] = char

	return nil
}
