package sm

import (
	"context"
	"errors"
)

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

type UsersRepository interface {
	Save(ctx context.Context, user User) error
	User(ctx context.Context, username string) (User, error)
}
