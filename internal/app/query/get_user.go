package query

import (
	"context"
	"log/slog"
	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type GetUser struct {
	Username string
}

type GetUserHandler decorator.QueryHandler[GetUser, User]

type getUserHandler struct {
	users sm.UsersRepository
}

func NewGetUserHandler(
	users sm.UsersRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetUserHandler {
	if users == nil {
		panic("users repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetUser, User](
		&getUserHandler{users},
		log, metricsClient,
	)
}

func (h *getUserHandler) Handle(ctx context.Context, q GetUser) (User, error) {
	user, err := h.users.User(ctx, q.Username)
	if err != nil {
		return User{}, err
	}

	return convertUserToApp(user), nil
}
