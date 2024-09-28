package mocks

import (
	"context"
	"sync"

	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

type mockUsersRepository struct {
	m map[string]sm.User
	sync.RWMutex
}

func NewMockUsersRepository() sm.UsersRepository {
	r := &mockUsersRepository{
		m: map[string]sm.User{},
	}

	_ = r.Save(nil, sm.MustNewUser("zhikhkirill", sm.Participant))

	return r
}

func (r *mockUsersRepository) Save(_ context.Context, user sm.User) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[user.Username]; ok {
		return sm.ErrUserAlreadyExists
	}

	r.m[user.Username] = user

	return nil
}

func (r *mockUsersRepository) User(_ context.Context, username string) (sm.User, error) {
	r.RLock()
	defer r.RUnlock()

	user, ok := r.m[username]
	if !ok {
		return sm.User{}, sm.ErrUserNotFound
	}

	return user, nil
}
