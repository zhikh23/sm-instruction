package telegram

import (
	"log/slog"

	"sm-instruction/internal/app"
	"sm-instruction/internal/common/logs"
)

type Port struct {
	app *app.Application
	log *slog.Logger
}

func NewTelegramPort(app *app.Application) *Port {
	log := logs.DefaultLogger()

	return &Port{
		app: app,
		log: log,
	}
}
