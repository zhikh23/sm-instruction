package sm

import (
	"github.com/zhikh23/sm-instruction/internal/common/commonerrs"
)

type User struct {
	Username string
	Role     Role
}

func (u User) IsZero() bool {
	return u == User{}
}

func NewUser(
	username string,
	role Role,
) (User, error) {
	if username == "" {
		return User{}, commonerrs.NewInvalidInputError("expected not empty username")
	}

	if role.IsZero() {
		return User{}, commonerrs.NewInvalidInputError("expected not empty role")
	}

	return User{
		Username: username,
		Role:     role,
	}, nil
}

func MustNewUser(
	username string,
	role Role,
) User {
	u, err := NewUser(username, role)
	if err != nil {
		panic(err)
	}
	return u
}

func UnmarshallUserFromDB(
	username string,
	role string,
) (User, error) {
	if username == "" {
		return User{}, commonerrs.NewInvalidInputError("expected not empty username")
	}

	if role == "" {
		return User{}, commonerrs.NewInvalidInputError("expected not empty role")
	}

	r, err := NewRoleFromString(role)
	if err != nil {
		return User{}, err
	}

	return User{
		Username: username,
		Role:     r,
	}, nil
}
