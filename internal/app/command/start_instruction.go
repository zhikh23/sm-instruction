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
	logs *slog.Logger,
	metricsClient decorator.MetricsClient,
) StartInstructionHandler {
	return decorator.ApplyCommandDecorators[StartInstruction](
		&startInstructionHandler{},
		logs,
		metricsClient,
	)
}

func (h *startInstructionHandler) Handle(ctx context.Context, cmd StartInstruction) error {
	return nil
}
