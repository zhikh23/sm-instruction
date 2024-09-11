package mocks

import (
	"context"
	"sync"

	"sm-instruction/internal/domain/sm"
)

type mockUsersRepository struct {
	m map[string]sm.User
	sync.RWMutex
}

func NewMockUsersRepository() sm.UsersRepository {
	return &mockUsersRepository{
		m: map[string]sm.User{},
	}
}

func (r *mockUsersRepository) Upsert(ctx context.Context, user sm.User) error {
	r.Lock()
	defer r.Unlock()

	r.m[user.Username] = user

	return nil
}

func (r *mockUsersRepository) User(ctx context.Context, username string) (sm.User, error) {
	r.RLock()
	defer r.RUnlock()

	user, ok := r.m[username]
	if !ok {
		return sm.User{}, sm.ErrUserNotFound
	}

	return user, nil
}
