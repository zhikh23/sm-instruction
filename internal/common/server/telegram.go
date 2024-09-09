package server

import (
	"github.com/vitaliy-ukiru/telebot-filter/dispatcher"
	"os"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/pkg/storage/memory"
)

func RunTelegramServer(setupFn func(m *fsm.Manager, dp fsm.Dispatcher)) {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		panic("TELEGRAM_TOKEN environment variable not set")
	}
	RunTelegramServerWithToken(token, setupFn)
}

func RunTelegramServerWithToken(token string, setupFn func(m *fsm.Manager, dp fsm.Dispatcher)) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		panic(err)
	}

	g := bot.Group()
	dp := dispatcher.NewDispatcher(g)

	m := fsm.New(memory.NewStorage())
	g.Use(m.WrapContext)

	setupFn(m, dp)

	bot.Start()
}
