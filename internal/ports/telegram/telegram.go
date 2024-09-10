package telegram

import (
	"context"
	"errors"
	"fmt"
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
	menuCommandState        = fsm.State("menuCommandState")

	bookChooseLocationState = fsm.State("bookChooseLocationState")
	bookChooseTimeState     = fsm.State("bookChooseTimeState")
)

const (
	toBookButton       = "Забронировать точку"
	toCancelBookButton = "Отменить бронирование точки"
	toCharacterProfile = "Перейти в профиль"
)

func (p *Port) menuReplyMarkup(char query.Character) *telebot.ReplyMarkup {
	opts := make([]string, 0)

	opts = append(opts, toCharacterProfile)
	if char.BookedLocationUUID == nil {
		opts = append(opts, toBookButton)
	} else {
		opts = append(opts, toCancelBookButton)
	}

	return &telebot.ReplyMarkup{
		ReplyKeyboard: makeReplyButtonsForMenu(opts),
	}
}

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
		fsmopt.On(toBookButton),
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

	dp.Dispatch(m.New(
		fsmopt.On(toCharacterProfile),
		fsmopt.OnStates(menuCommandState),
		fsmopt.Do(p.profile),
	))

	dp.Dispatch(m.New(
		fsmopt.On(toCancelBookButton),
		fsmopt.OnStates(menuCommandState),
		fsmopt.Do(p.cancelBooking),
	))
}

func (p *Port) cancel(c telebot.Context, state fsm.Context) error {
	_ = state.Finish(context.Background(), c.Data() != "")
	return c.Send("Отменено.")
}

func (p *Port) start(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	char, err := p.app.Queries.GetCharacter.Handle(ctx, query.GetCharacter{ChatID: c.Chat().ID})
	if errors.Is(err, sm.ErrCharacterNotFound) {
		err = state.SetState(ctx, startReadGroupNameState)
		if err != nil {
			return err
		}

		return c.Send(
			"Привет, участник «СМ. Инструкция по выживанию»!\n" +
				"Для начала, напиши пожалуйста свою группу, например: СМ1-11Б.",
		)
	} else if err != nil {
		return err
	}

	return c.Send("Что нужно сделать?", p.menuReplyMarkup(char))
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

	char, err := p.app.Queries.GetCharacter.Handle(ctx, query.GetCharacter{ChatID: c.Chat().ID})
	if err != nil {
		return err
	}

	err = state.SetState(ctx, menuCommandState)
	if err != nil {
		return err
	}

	return c.Send(
		"Хорошо, теперь ты можешь приступать к инструкции!\n"+
			"У тебя есть 4 часа на прохождение всех точек, поэтому поторопись :)\n"+
			"Полный список команд ты можешь узнать при помощи /help.", p.menuReplyMarkup(char),
	)
}

func (p *Port) menu(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	err := state.SetState(ctx, menuCommandState)
	if err != nil {
		return err
	}

	return c.Send(fmt.Sprintf(
		"Доступные команды:\n"+
			"'%s' - забронировать точку;\n"+
			"'%s' посмотреть профиль персонажа",
		toBookButton, toCharacterProfile,
	), &telebot.ReplyMarkup{
		ReplyKeyboard: makeReplyButtonsForMenu([]string{toBookButton, toCharacterProfile}),
	})
}

func (p *Port) book(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	char, err := p.app.Queries.GetCharacter.Handle(ctx, query.GetCharacter{ChatID: c.Chat().ID})
	if err != nil {
		return err
	}

	if char.BookedLocationUUID != nil {
		return c.Send("")
	}

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

	char, err := p.app.Queries.GetCharacter.Handle(ctx, query.GetCharacter{ChatID: c.Chat().ID})
	if err != nil {
		return err
	}

	err = state.SetState(ctx, menuCommandState)
	if err != nil {
		return err
	}

	return c.Send(
		"Успешно забронировано! "+
			"Местоположение точки и время бронирования можешь посмотреть в своем профиле.",
		p.menuReplyMarkup(char),
	)
}

func (p *Port) cancelBooking(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	err := p.app.Commands.CancelBooking.Handle(ctx, command.CancelBooking{ChatID: c.Chat().ID})
	if err != nil {
		return err
	}

	char, err := p.app.Queries.GetCharacter.Handle(ctx, query.GetCharacter{ChatID: c.Chat().ID})
	if err != nil {
		return err
	}

	err = state.SetState(ctx, menuCommandState)
	if err != nil {
		return err
	}

	return c.Send("Бронирование отменено.", p.menuReplyMarkup(char))
}

func (p *Port) profile(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	char, err := p.app.Queries.GetCharacter.Handle(ctx, query.GetCharacter{ChatID: c.Chat().ID})
	if err != nil {
		return err
	}

	msg := "Данные о твоём персонаже:\n"
	msg += fmt.Sprintf("Учебная группа: %s\n", char.GroupName)
	msg += fmt.Sprintf("Начал в: %s\n", char.StartedAt.Format(timeFormat))
	msg += fmt.Sprintf("Конец Инструкции: %s\n", char.FinishAt.Format(timeFormat))

	if char.BookedLocationUUID != nil {
		uuid := *char.BookedLocationUUID
		loc, err := p.app.Queries.GetLocation.Handle(ctx, query.GetLocation{UUID: uuid})
		if err != nil {
			return err
		}
		msg += fmt.Sprintf("Забронировано: '%s' до %s\n", loc.Name, char.BookedLocationTo.Format(timeFormat))
		msg += fmt.Sprintf("Местонахождение точки: %s\n", loc.Where)
	}

	err = state.SetState(ctx, menuCommandState)
	if err != nil {
		return err
	}

	return c.Send(msg, p.menuReplyMarkup(char))
}

func makeReplyButtonsForMenu(opts []string) [][]telebot.ReplyButton {
	buttons := make([]telebot.ReplyButton, len(opts))

	for i, opt := range opts {
		buttons[i] = telebot.ReplyButton{Text: opt}
	}

	return composeButtons(buttons, 1)
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
