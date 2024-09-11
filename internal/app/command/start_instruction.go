package command

import (
	"context"
	"log/slog"

	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/domain/sm"
)

type StartInstruction struct {
	ChatID    int64
	Username  string
	GroupName string
}

type StartInstructionHandler decorator.CommandHandler[StartInstruction]

type startInstructionHandler struct {
	users sm.UsersRepository
	chars sm.CharactersRepository
}

func NewStartInstructionHandler(
	users sm.UsersRepository,
	chars sm.CharactersRepository,
	logs *slog.Logger,
	metricsClient decorator.MetricsClient,
) StartInstructionHandler {
	return decorator.ApplyCommandDecorators[StartInstruction](
		&startInstructionHandler{users: users, chars: chars},
		logs,
		metricsClient,
	)
}

func (h *startInstructionHandler) Handle(ctx context.Context, cmd StartInstruction) error {
	user, err := sm.NewUser(cmd.Username, sm.Participant)
	if err != nil {
		return err
	}

	char, err := sm.NewCharacter(cmd.Username, cmd.GroupName)
	if err != nil {
		return err
	}

	if err = char.Start(); err != nil {
		return err
	}

	if err = h.users.Upsert(ctx, user); err != nil {
		return err
	}

	return h.chars.Save(ctx, char)
}
