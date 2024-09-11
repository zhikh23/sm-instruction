package sm

import (
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type UsersRepository interface {
	Upsert(ctx context.Context, user User) error
	User(ctx context.Context, username string) (User, error)
}
