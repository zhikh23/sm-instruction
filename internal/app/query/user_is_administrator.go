package query

import (
	"context"
	"errors"
	"log/slog"

	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type UserIsAdministrator struct {
	Username string
}

type UserIsAdministratorHandler decorator.QueryHandler[UserIsAdministrator, bool]

type userIsAdministratorHandler struct {
	users sm.UsersRepository
}

func NewUserIsAdministratorHandler(
	users sm.UsersRepository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) UserIsAdministratorHandler {
	if users == nil {
		panic("users repository is nil")
	}

	return decorator.ApplyQueryDecorators[UserIsAdministrator, bool](
		&userIsAdministratorHandler{users: users},
		log,
		metricsClient,
	)
}

func (h *userIsAdministratorHandler) Handle(ctx context.Context, query UserIsAdministrator) (bool, error) {
	user, err := h.users.User(ctx, query.Username)
	if errors.Is(err, sm.ErrUserNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return user.Role == sm.Administrator, nil
}
