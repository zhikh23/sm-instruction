package command

import (
	"context"
	"log/slog"

	"github.com/zhikh23/sm-instruction/internal/common/decorator"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

type StartInstruction struct {
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
	return h.chars.Update(ctx, cmd.GroupName, func(innerCtx context.Context, char *sm.Character) error {
		return char.Start()
	})
}
