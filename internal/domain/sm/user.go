package sm

import "sm-instruction/internal/common/commonerrs"

type User struct {
	ChatID   int64
	Username string
}

func (u User) IsZero() bool {
	return u == User{}
}

func NewUser(
	chatID int64,
	username string,
) (User, error) {
	if chatID == 0 {
		return User{}, commonerrs.NewInvalidInputError("expected not empty chat ID")
	}

	if username == "" {
		return User{}, commonerrs.NewInvalidInputError("expected not empty username")
	}

	return User{
		ChatID:   chatID,
		Username: username,
	}, nil
}

func MustNewUser(
	chatID int64,
	username string,
) User {
	u, err := NewUser(chatID, username)
	if err != nil {
		panic(err)
	}
	return u
}
