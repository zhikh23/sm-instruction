package telegram

import (
	"log/slog"

	"github.com/zhikh23/sm-instruction/internal/app"
	"github.com/zhikh23/sm-instruction/internal/common/logs"
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
