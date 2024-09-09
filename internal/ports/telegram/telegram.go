package telegram

import (
	"context"
	"errors"
	"log/slog"
	"math"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/fsmopt"
	"gopkg.in/telebot.v3"

	"sm-instruction/internal/app"
	"sm-instruction/internal/app/command"
	"sm-instruction/internal/app/query"
	"sm-instruction/internal/common/commonerrs"
	"sm-instruction/internal/common/logs"
	"sm-instruction/internal/domain/sm"
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

const (
	startReadGroupNameState = fsm.State("startReadGroupNameState")
	startedState            = fsm.State("startedState")
	menuCommandState        = fsm.State("menuCommandState")

	bookChooseLocationState = fsm.State("bookChooseLocationState")
	bookChooseTimeState     = fsm.State("bookChooseTimeState")
)

func (p *Port) RegisterFSMManager(m *fsm.Manager, dp fsm.Dispatcher) {
	dp.Dispatch(m.New(
		fsmopt.On("/cancel"),
		fsmopt.OnStates(fsm.AnyState),
		fsmopt.Do(p.cancel),
	))

	dp.Dispatch(m.New(
		fsmopt.On("/start"),
		fsmopt.OnStates(fsm.AnyState),
		fsmopt.Do(p.start),
	))

	dp.Dispatch(m.New(
		fsmopt.On(telebot.OnText),
		fsmopt.OnStates(startReadGroupNameState),
		fsmopt.Do(p.startReadGroupName),
	))

	dp.Dispatch(m.New(
		fsmopt.On("/menu"),
		fsmopt.OnStates(startedState),
		fsmopt.Do(p.menu),
	))

	dp.Dispatch(m.New(
		fsmopt.On("/book"),
		fsmopt.OnStates(menuCommandState),
		fsmopt.Do(p.book),
	))

	dp.Dispatch(m.New(
		fsmopt.On(telebot.OnText),
		fsmopt.OnStates(bookChooseLocationState),
		fsmopt.Do(p.bookChooseLocation),
	))

	dp.Dispatch(m.New(
		fsmopt.On(telebot.OnText),
		fsmopt.OnStates(bookChooseTimeState),
		fsmopt.Do(p.bookChooseTime),
	))
}

func (p *Port) cancel(c telebot.Context, state fsm.Context) error {
	_ = state.Finish(context.Background(), c.Data() != "")
	return c.Send("Отменено.")
}

func (p *Port) start(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	err := c.Send(
		"Привет, участник «СМ. Инструкция по выживанию»!\n" +
			"Для начала, напиши пожалуйста свою группу, например: СМ1-11Б.",
	)
	if err != nil {
		return err
	}

	return state.SetState(ctx, startReadGroupNameState)
}

func (p *Port) startReadGroupName(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	groupName := c.Message().Text

	err := p.app.Commands.StartInstruction.Handle(ctx, command.StartInstruction{
		ChatID:    c.Chat().ID,
		Username:  c.Chat().Username,
		GroupName: groupName,
	})
	if errors.As(err, &commonerrs.InvalidInputError{}) {
		return c.Send("Некорректное название группы, попробуй ещё раз!")
	} else if err != nil {
		return err
	}

	err = state.SetState(ctx, startedState)
	if err != nil {
		return err
	}

	return c.Send(
		"Хорошо, теперь ты можешь приступать к инструкции!\n" +
			"У тебя есть 4 часа на прохождение всех точек, поэтому поторопись :)\n" +
			"Полный список команд ты можешь узнать в /menu.",
	)
}

func (p *Port) menu(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	err := state.SetState(ctx, menuCommandState)
	if err != nil {
		return err
	}

	return c.Send(
		"Доступные команды:\n" +
			"/book - забронировать точку;\n" +
			"/cancel - отменить операцию.\n",
	)
}

func (p *Port) book(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	locs, err := p.app.Queries.GetAllLocations.Handle(ctx, query.GetAllLocations{})
	if err != nil {
		return err
	}

	err = state.SetState(ctx, bookChooseLocationState)
	if err != nil {
		return err
	}

	return c.Send("Выбери точку из доступных:", &telebot.ReplyMarkup{
		ReplyKeyboard: makeReplyButtonsForLocations(locs),
	})
}

func (p *Port) bookChooseLocation(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	chatID := c.Chat().ID
	locationName := c.Message().Text

	loc, err := p.app.Queries.GetLocationByName.Handle(ctx, query.GetLocationByName{Name: locationName})
	if errors.Is(err, sm.ErrLocationNotFound) {
		return c.Send("Такой точки не существует :( Выбери, пожалуйста, точку из списка.")
	} else if err != nil {
		return err
	}

	err = state.Update(ctx, "locationUUID", loc.UUID)

	available, err := p.app.Queries.GetAvailableIntervals.Handle(ctx, query.GetAvailableIntervals{
		ChatID:       chatID,
		LocationUUID: loc.UUID,
	})
	if err != nil {
		return err
	}

	if len(available) == 0 {
		return c.Send("Ой, у точки не осталось времени, которое можно забронировать :( Попробуй выбрать другую точку.")
	}

	err = state.SetState(ctx, bookChooseTimeState)
	if err != nil {
		return err
	}

	return c.Send(
		"Выбери доступный диапазон времени для бронирования точки:", &telebot.ReplyMarkup{
			ReplyKeyboard: makeReplyButtonsFromTimes(available),
		},
	)
}

const timeFormat = "15:04"

func (p *Port) bookChooseTime(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	t, err := closestTimeFromString(c.Message().Text)
	if err != nil {
		return c.Send("Пожалуйста, выбери одно из доступных времени бронирования.")
	}

	var locationUUID string
	err = state.Data(ctx, "locationUUID", &locationUUID)
	if err != nil {
		return err
	}

	err = p.app.Commands.BookLocation.Handle(ctx, command.BookLocation{
		LocationUUID: locationUUID,
		ChatID:       c.Chat().ID,
		Time:         t,
	})
	switch {
	case errors.As(err, &commonerrs.InvalidInputError{}):
		return c.Send("Пожалуйста, выбери корректное время бронирования.")
	case errors.Is(err, sm.ErrLocationAlreadyBooked):
		return c.Send("К сожалению, данное время уже забронировано. Пожалуйста, выбери другое время.")
	case errors.Is(err, sm.ErrCharacterBookingIsTooLate):
		return c.Send("Невозможно забронировать точку позже, чем 4 часа после начала Инструкции. Пожалуйста, выбери другое время.")
	case errors.Is(err, sm.ErrCharacterBookingIsTooClose):
		return c.Send("Невозможно забронировать точку позже чем за 5 минут до старта. Пожалуйста, выбери другое время.")
	case errors.Is(err, sm.ErrCharacterAlreadyFinished):
		return c.Send("Невозможно забронировать точку, т.к. ты уже завершил Инструкцию.")
	case err != nil:
		return err
	}

	err = state.SetState(ctx, menuCommandState)
	if err != nil {
		return err
	}

	return c.Send(
		"Успешно забронировано! "+
			"Местоположение точки и время бронирования можешь посмотреть в своем профиле (/profile)",
		telebot.ReplyMarkup{
			RemoveKeyboard: true,
		},
	)
}

func makeReplyButtonsForLocations(locs []query.Location) [][]telebot.ReplyButton {
	buttons := make([]telebot.ReplyButton, len(locs))

	for i, loc := range locs {
		buttons[i] = telebot.ReplyButton{Text: loc.Name}
	}

	return composeButtons(buttons, 2)
}

func makeReplyButtonsFromTimes(times []time.Time) [][]telebot.ReplyButton {
	buttons := make([]telebot.ReplyButton, len(times))

	for i, t := range times {
		buttons[i] = telebot.ReplyButton{
			Text: t.Format(timeFormat),
		}
	}

	return composeButtons(buttons, 3)
}

func closestTimeFromString(str string) (time.Time, error) {
	t, err := time.Parse(timeFormat, str)
	if err != nil {
		return time.Time{}, err
	}
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, time.Local), nil
}

func composeButtons(buttons []telebot.ReplyButton, width int) [][]telebot.ReplyButton {
	height := int(math.Ceil(float64(len(buttons)) / float64(width)))

	rows := make([][]telebot.ReplyButton, height)
	y := 0
	for ; y < height-1; y++ {
		rows[y] = buttons[y*width : (y+1)*width]
	}
	rows[y] = buttons[y*width:]

	return rows
}
