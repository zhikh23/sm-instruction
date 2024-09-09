package mocks

import (
	"context"
	"sync"

	"sm-instruction/internal/domain/sm"
)

type mockCharactersRepository struct {
	m map[int64]sm.Character
	sync.RWMutex
}

func NewMockCharactersRepository() sm.CharactersRepository {
	return &mockCharactersRepository{
		m: make(map[int64]sm.Character),
	}
}

func (r *mockCharactersRepository) Save(_ context.Context, char *sm.Character) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[char.Head.ChatID]; ok {
		return sm.ErrCharacterAlreadyExists
	}

	r.m[char.Head.ChatID] = *char

	return nil
}

func (r *mockCharactersRepository) Character(_ context.Context, chatID int64) (*sm.Character, error) {
	r.RLock()
	defer r.RUnlock()

	char, ok := r.m[chatID]
	if !ok {
		return nil, sm.ErrCharacterNotFound
	}

	return &char, nil
}

func (r *mockCharactersRepository) Update(
	ctx context.Context,
	chatID int64,
	updateFn func(innerCtx context.Context, char *sm.Character) error,
) error {
	r.Lock()
	defer r.Unlock()

	char, ok := r.m[chatID]
	if !ok {
		return sm.ErrCharacterNotFound
	}

	err := updateFn(ctx, &char)
	if err != nil {
		return err
	}

	r.m[chatID] = char

	return nil
}
