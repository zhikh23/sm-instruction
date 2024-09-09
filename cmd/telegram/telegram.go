package main

import (
	"github.com/vitaliy-ukiru/fsm-telebot/v2"

	"sm-instruction/internal/common/server"
	"sm-instruction/internal/ports/telegram"
	"sm-instruction/internal/service"
)

func main() {
	app := service.NewMockedApplication()

	server.RunTelegramServer(func(m *fsm.Manager, dp fsm.Dispatcher) {
		telegram.NewTelegramPort(app).RegisterFSMManager(m, dp)
	})
}
