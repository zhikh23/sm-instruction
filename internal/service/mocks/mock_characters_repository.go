package mocks

import (
	"context"
	"sync"

	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

type mockCharactersRepository struct {
	m map[string]sm.Character
	sync.RWMutex
}

func NewMockCharactersRepository() sm.CharactersRepository {
	r := &mockCharactersRepository{
		m: make(map[string]sm.Character),
	}

	_ = r.Save(nil, sm.MustNewCharacter(
		"СМ1-11Б",
		"zhikhkirill",
		[]*sm.Slot{
			sm.MustNewSlot(timeWithMinutes(0), timeWithMinutes(15)),
			sm.MustNewSlot(timeWithMinutes(15), timeWithMinutes(30)),
		},
	))

	return r
}

func (r *mockCharactersRepository) Save(_ context.Context, char *sm.Character) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[char.GroupName]; ok {
		return sm.ErrCharacterAlreadyExists
	}

	r.m[char.GroupName] = *char

	return nil
}

func (r *mockCharactersRepository) Character(_ context.Context, groupName string) (*sm.Character, error) {
	r.RLock()
	defer r.RUnlock()

	char, ok := r.m[groupName]
	if !ok {
		return nil, sm.ErrCharacterNotFound
	}

	return &char, nil
}

func (r *mockCharactersRepository) CharacterByUsername(_ context.Context, username string) (*sm.Character, error) {
	r.RLock()
	defer r.RUnlock()

	for _, char := range r.m {
		if char.Username == username {
			return &char, nil
		}
	}

	return nil, sm.ErrCharacterNotFound
}

func (r *mockCharactersRepository) Update(
	ctx context.Context,
	groupName string,
	updateFn func(innerCtx context.Context, char *sm.Character) error,
) error {
	r.Lock()
	defer r.Unlock()

	char, ok := r.m[groupName]
	if !ok {
		return sm.ErrCharacterNotFound
	}

	err := updateFn(ctx, &char)
	if err != nil {
		return err
	}

	r.m[groupName] = char

	return nil
}
