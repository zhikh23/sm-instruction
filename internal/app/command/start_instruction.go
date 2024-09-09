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
	chars sm.CharactersRepository
}

func NewStartInstructionHandler(
	chars sm.CharactersRepository,
	logs *slog.Logger,
	metricsClient decorator.MetricsClient,
) StartInstructionHandler {
	return decorator.ApplyCommandDecorators[StartInstruction](
		&startInstructionHandler{chars: chars},
		logs,
		metricsClient,
	)
}

func (h *startInstructionHandler) Handle(ctx context.Context, cmd StartInstruction) error {
	user, err := sm.NewUser(cmd.ChatID, cmd.Username)
	if err != nil {
		return err
	}

	char, err := sm.NewCharacter(user, cmd.GroupName)
	if err != nil {
		return err
	}

	if err = char.Start(); err != nil {
		return err
	}

	return h.chars.Save(ctx, char)
}
