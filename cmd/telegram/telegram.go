package main

import (
	"log"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"

	"github.com/zhikh23/sm-instruction/internal/common/server"
	"github.com/zhikh23/sm-instruction/internal/ports/telegram"
	"github.com/zhikh23/sm-instruction/internal/service"
)

func main() {
	app, closeFn := service.NewApplication()
	defer func() {
		err := closeFn()
		if err != nil {
			log.Fatal(err)
		}
	}()

	server.RunTelegramServer(func(m *fsm.Manager, dp fsm.Dispatcher) {
		telegram.NewTelegramPort(app).RegisterFSMManager(m, dp)
	})
}
