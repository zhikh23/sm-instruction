package telegram

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strconv"
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
	participantMenuState    = fsm.State("participantMenuState")
	bookChooseLocationState = fsm.State("bookChooseLocationState")
	bookApproveLocation     = fsm.State("bookApproveLocation")
	bookChooseTimeState     = fsm.State("bookChooseTimeState")

	adminMenuState = fsm.State("adminMenuState")

	awardEnterGroupNameState = fsm.State("awardEnterGroupNameState")
	awardEnterSkillTypeState = fsm.State("awardEnterSkillTypeState")
	awardEnterPointsState    = fsm.State("awardEnterPointsState")
)

const (
	approveButton = "Подтвердить"
	cancelButton  = "Отменить"

	toBookButton       = "Забронировать точку"
	toCancelBookButton = "Отменить бронирование точки"
	toCharacterProfile = "Перейти в профиль"

	toShowBookings   = "Показать забронированные точки"
	toAwardCharacter = "Начислить участнику баллов"
)

func (p *Port) participantMenuReplyMarkup(char query.Character) *telebot.ReplyMarkup {
	opts := make([]string, 0)

	opts = append(opts, toCharacterProfile)
	if char.BookedTime == nil {
		opts = append(opts, toBookButton)
	} else if char.BookedTime.CanBeRemoved {
		opts = append(opts, toCancelBookButton)
	}

	return &telebot.ReplyMarkup{
		ReplyKeyboard: makeReplyButtonsForMenu(opts),
	}
}

func (p *Port) locationAdminMenuReplyMarkup() *telebot.ReplyMarkup {
	opts := make([]string, 0)

	opts = append(opts, toShowBookings)
	opts = append(opts, toAwardCharacter)

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
		fsmopt.On("/admin"),
		fsmopt.OnStates(fsm.AnyState),
		fsmopt.Do(p.adminMenu),
	))

	dp.Dispatch(m.New(
		fsmopt.On(telebot.OnText),
		fsmopt.OnStates(startReadGroupNameState),
		fsmopt.Do(p.startReadGroupName),
	))

	dp.Dispatch(m.New(
		fsmopt.On(toBookButton),
		fsmopt.OnStates(participantMenuState),
		fsmopt.Do(p.book),
	))

	dp.Dispatch(m.New(
		fsmopt.On(telebot.OnText),
		fsmopt.OnStates(bookChooseLocationState),
		fsmopt.Do(p.bookLocationDescribe),
	))

	dp.Dispatch(m.New(
		fsmopt.On(approveButton),
		fsmopt.OnStates(bookApproveLocation),
		fsmopt.Do(p.bookChooseLocationTime),
	))

	dp.Dispatch(m.New(
		fsmopt.On(cancelButton),
		fsmopt.OnStates(bookApproveLocation),
		fsmopt.Do(p.book),
	))

	dp.Dispatch(m.New(
		fsmopt.On(toCharacterProfile),
		fsmopt.OnStates(participantMenuState),
		fsmopt.Do(p.profile),
	))

	dp.Dispatch(m.New(
		fsmopt.On(telebot.OnText),
		fsmopt.OnStates(bookChooseTimeState),
		fsmopt.Do(p.bookChooseTime),
	))

	dp.Dispatch(m.New(
		fsmopt.On(toCancelBookButton),
		fsmopt.OnStates(participantMenuState),
		fsmopt.Do(p.cancelBooking),
	))

	dp.Dispatch(m.New(
		fsmopt.On(toShowBookings),
		fsmopt.OnStates(adminMenuState),
		fsmopt.Do(p.showBookings),
	))

	dp.Dispatch(m.New(
		fsmopt.On(toAwardCharacter),
		fsmopt.OnStates(adminMenuState),
		fsmopt.Do(p.awardEnterCharacterGroup),
	))

	dp.Dispatch(m.New(
		fsmopt.On(telebot.OnText),
		fsmopt.OnStates(awardEnterGroupNameState),
		fsmopt.Do(p.acceptCharacterGroupName),
	))

	dp.Dispatch(m.New(
		fsmopt.On(telebot.OnText),
		fsmopt.OnStates(awardEnterSkillTypeState),
		fsmopt.Do(p.acceptSkillType),
	))

	dp.Dispatch(m.New(
		fsmopt.On(telebot.OnText),
		fsmopt.OnStates(awardEnterPointsState),
		fsmopt.Do(p.acceptPoints),
	))
}

func (p *Port) cancel(c telebot.Context, state fsm.Context) error {
	_ = state.Finish(context.Background(), c.Data() != "")
	return c.Send("Отменено.")
}

func (p *Port) start(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	isAdmin, err := p.app.Queries.UserIsAdministrator.Handle(ctx, query.UserIsAdministrator{Username: c.Chat().Username})
	if err != nil {
		return err
	}

	if isAdmin {
		return p.adminMenu(c, state)
	}

	isStarted, err := p.app.Queries.CharacterIsStarted.Handle(ctx, query.CharacterIsStarted{Username: c.Chat().Username})
	if err != nil {
		return err
	}

	if !isStarted {
		err = state.SetState(ctx, startReadGroupNameState)
		if err != nil {
			return err
		}

		return c.Send(
			"Привет, участник «СМ. Инструкция по выживанию»!\n" +
				"Для начала, напиши пожалуйста свою группу, например: СМ1-11Б.",
		)
	}

	char, err := p.app.Queries.GetCharacter.Handle(ctx, query.GetCharacter{Username: c.Chat().Username})
	if err != nil {
		return err
	}

	if err = state.SetState(ctx, participantMenuState); err != nil {
		return err
	}

	return c.Send("Что нужно сделать?", p.participantMenuReplyMarkup(char))
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

	char, err := p.app.Queries.GetCharacter.Handle(ctx, query.GetCharacter{Username: c.Chat().Username})
	if err != nil {
		return err
	}

	if err = state.SetState(ctx, participantMenuState); err != nil {
		return err
	}

	return c.Send(
		"Хорошо, теперь ты можешь приступать к инструкции!\n"+
			"У тебя есть 4 часа на прохождение всех точек, поэтому поторопись :)\n"+
			"Полный список команд ты можешь узнать при помощи /help.", p.participantMenuReplyMarkup(char),
	)
}

func (p *Port) participantMenu(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	if err := state.SetState(ctx, participantMenuState); err != nil {
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

func (p *Port) adminMenu(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	act, err := p.app.Queries.GetActivityByAdmin.Handle(ctx, query.GetActivityByAdmin{Username: c.Chat().Username})
	if err != nil {
		return err
	}

	if err = state.Update(ctx, "activityUUID", act.UUID); err != nil {
		return err
	}

	if err = state.SetState(ctx, adminMenuState); err != nil {
		return err
	}

	return c.Send("Панель управления администратора. Выбери действие:", p.locationAdminMenuReplyMarkup())
}

func (p *Port) book(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	locs, err := p.app.Queries.GetAllLocations.Handle(ctx, query.GetAllLocations{})
	if err != nil {
		return err
	}

	if err = state.SetState(ctx, bookChooseLocationState); err != nil {
		return err
	}

	return c.Send("Выбери точку из доступных:", &telebot.ReplyMarkup{
		ReplyKeyboard: makeReplyButtonsForLocations(locs),
	})
}

func (p *Port) bookLocationDescribe(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	locationName := c.Message().Text

	loc, err := p.app.Queries.GetLocationByName.Handle(ctx, query.GetLocationByName{Name: locationName})
	if errors.Is(err, sm.ErrActivityNotFound) || errors.Is(err, sm.ErrActivityHasNotLocation) {
		return c.Send("Такой точки не существует :( Выбери, пожалуйста, точку из списка.")
	} else if err != nil {
		return err
	}

	if err = state.Update(ctx, "locationUUID", loc.UUID); err != nil {
		return err
	}

	if err = state.SetState(ctx, bookApproveLocation); err != nil {
		return err
	}

	return c.Send(
		"Описание точки:\n\n"+loc.Description+"\n Подтверждаешь бронирование?",
		&telebot.ReplyMarkup{
			ReplyKeyboard: makeReplyButtonsApprove(),
		},
	)
}

func (p *Port) bookChooseLocationTime(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	var locationUUID string
	err := state.Data(ctx, "locationUUID", &locationUUID)
	if err != nil {
		return err
	}

	available, err := p.app.Queries.GetAvailableIntervals.Handle(ctx, query.GetAvailableIntervals{
		Username:     c.Chat().Username,
		ActivityUUID: locationUUID,
	})
	if err != nil {
		return err
	}

	if len(available) == 0 {
		return c.Send("Ой, у точки не осталось времени, которое можно забронировать :( Попробуй выбрать другую точку.")
	}

	if err = state.SetState(ctx, bookChooseTimeState); err != nil {
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
		Username:     c.Chat().Username,
		Time:         t,
	})
	switch {
	case errors.As(err, &commonerrs.InvalidInputError{}):
		return c.Send("Пожалуйста, выбери корректное время бронирования.")
	case errors.Is(err, sm.ErrLocationAlreadyBooked):
		return c.Send("К сожалению, данное время уже забронировано. Пожалуйста, выбери другое время.")
	case errors.Is(err, sm.ErrCharacterBookingAfterFinish):
		return c.Send("Невозможно забронировать точку позже, чем 4 часа после начала Инструкции. Пожалуйста, выбери другое время.")
	case errors.Is(err, sm.ErrCharacterBookingIsTooClose):
		return c.Send("Невозможно забронировать точку позже чем за 5 минут до старта. Пожалуйста, выбери другое время.")
	case err != nil:
		return err
	}

	char, err := p.app.Queries.GetCharacter.Handle(ctx, query.GetCharacter{Username: c.Chat().Username})
	if err != nil {
		return err
	}

	if err = state.SetState(ctx, participantMenuState); err != nil {
		return err
	}

	return c.Send(
		"Успешно забронировано! "+
			"Местоположение точки и время бронирования можешь посмотреть в своем профиле.",
		p.participantMenuReplyMarkup(char),
	)
}

func (p *Port) cancelBooking(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	err := p.app.Commands.RemoveBooking.Handle(ctx, command.RemoveBooking{Username: c.Chat().Username})
	if err != nil {
		return err
	}

	char, err := p.app.Queries.GetCharacter.Handle(ctx, query.GetCharacter{Username: c.Chat().Username})
	if err != nil {
		return err
	}

	if err = state.SetState(ctx, participantMenuState); err != nil {
		return err
	}

	return c.Send("Бронирование отменено.", p.participantMenuReplyMarkup(char))
}

func (p *Port) profile(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	char, err := p.app.Queries.GetCharacter.Handle(ctx, query.GetCharacter{Username: c.Chat().Username})
	if err != nil {
		return err
	}

	msg := "Данные о твоём персонаже:\n"
	msg += fmt.Sprintf("Учебная группа: %s\n", char.GroupName)
	msg += fmt.Sprintf("Начал в: %s\n", char.StartedAt.Format(timeFormat))
	msg += fmt.Sprintf("Конец Инструкции: %s\n", char.FinishAt.Format(timeFormat))
	msg += "\n"

	if char.BookedTime != nil {
		uuid := char.BookedTime.ActivityUUID
		loc, err := p.app.Queries.GetLocation.Handle(ctx, query.GetLocation{UUID: uuid})
		if err != nil {
			return err
		}
		msg += fmt.Sprintf(
			"Забронировано: %q на промежуток %s - %s\n",
			loc.Name,
			char.BookedTime.Start.Format(timeFormat),
			char.BookedTime.Finish.Format(timeFormat),
		)
		msg += fmt.Sprintf("Местонахождение точки: %s\n", loc.Where)
		msg += "\n"
	}

	msg += fmt.Sprintf("Твои навыки:\n")
	for skill, value := range char.Skills {
		msg += fmt.Sprintf("- %s - %d\n", skill, value)
	}

	if err = state.SetState(ctx, participantMenuState); err != nil {
		return err
	}

	return c.Send(msg, p.participantMenuReplyMarkup(char))
}

func (p *Port) showBookings(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	act, err := p.app.Queries.GetActivityByAdmin.Handle(ctx, query.GetActivityByAdmin{Username: c.Chat().Username})
	if err != nil {
		return err
	}

	loc, err := p.app.Queries.GetLocation.Handle(ctx, query.GetLocation{UUID: act.UUID})
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Точка '%s'\n", loc.Name)
	if len(loc.BookedTimes) == 0 {
		msg += "Не имеет бронирований."
	}
	for _, t := range loc.BookedTimes {
		msg += fmt.Sprintf(
			"<code>%s - %s</code> - @%s\n",
			t.Start.Format(timeFormat), t.Finish.Format(timeFormat), t.Username,
		)
	}

	if err = state.SetState(ctx, adminMenuState); err != nil {
		return err
	}

	return c.Send(msg, p.locationAdminMenuReplyMarkup(), telebot.ModeHTML)
}

func (p *Port) awardEnterCharacterGroup(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	if err := state.SetState(ctx, awardEnterGroupNameState); err != nil {
		return err
	}

	return c.Send(
		"Введите учебную группу участников: (в формате СМ1-11Б)",
		&telebot.ReplyMarkup{RemoveKeyboard: true},
	)
}

func (p *Port) acceptCharacterGroupName(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	groupName := c.Message().Text

	char, err := p.app.Queries.GetCharacterByGroup.Handle(ctx, query.GetCharacterByGroup{GroupName: groupName})
	if errors.Is(err, sm.ErrCharacterNotFound) {
		return c.Send("Упс, такой группы не существует :( Попробуй ещё раз.")
	}

	if err := state.Update(ctx, "username", char.Username); err != nil {
		return err
	}

	return p.awardEnterCharacterSkill(c, state)
}

func (p *Port) awardEnterCharacterSkill(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	var activityUUID string
	if err := state.Data(ctx, "activityUUID", &activityUUID); err != nil {
		return err
	}

	act, err := p.app.Queries.GetActivity.Handle(ctx, query.GetActivity{ActivityUUID: activityUUID})
	if err != nil {
		return err
	}

	if err = state.SetState(ctx, awardEnterSkillTypeState); err != nil {
		return err
	}

	return c.Send("Выбери тип навыков для начисления очков.",
		&telebot.ReplyMarkup{
			ReplyKeyboard: makeReplyButtonsForSkillTypes(act.Skills),
		},
	)
}

func (p *Port) acceptSkillType(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	skillType := c.Message().Text

	if _, err := sm.NewSkillTypeFromString(skillType); err != nil {
		return c.Send("Выбери один из доступных навыков для начисления очков.")
	}

	if err := state.Update(ctx, "skillType", skillType); err != nil {
		return err
	}

	return p.awardEnterPoints(c, state)
}

func (p *Port) awardEnterPoints(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	var activityUUID string
	if err := state.Data(ctx, "activityUUID", &activityUUID); err != nil {
		return err
	}

	act, err := p.app.Queries.GetActivity.Handle(ctx, query.GetActivity{ActivityUUID: activityUUID})
	if err != nil {
		return err
	}

	if err = state.SetState(ctx, awardEnterPointsState); err != nil {
		return err
	}

	return c.Send("Выбери количество очков для начисления участнику.",
		&telebot.ReplyMarkup{
			ReplyKeyboard: makeReplyButtonsForRangeNumbers(1, act.MaxPoints),
		},
	)
}

func (p *Port) acceptPoints(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	pointsStr := c.Message().Text
	points, err := strconv.Atoi(pointsStr)
	if err != nil {
		return c.Send("Пожалуйста, введи корректное число.")
	}

	if err = state.Update(ctx, "points", points); err != nil {
		return err
	}

	return p.awardCharacter(c, state)
}

func (p *Port) awardCharacter(c telebot.Context, state fsm.Context) error {
	ctx := context.Background()

	var activityUUID string
	if err := state.Data(ctx, "activityUUID", &activityUUID); err != nil {
		return err
	}

	var username string
	if err := state.Data(ctx, "username", &username); err != nil {
		return err
	}

	var skillType string
	if err := state.Data(ctx, "skillType", &skillType); err != nil {
		return err
	}

	var points int
	if err := state.Data(ctx, "points", &points); err != nil {
		return err
	}

	err := p.app.Commands.AwardCharacter.Handle(ctx, command.AwardCharacter{
		Username:     username,
		ActivityUUID: activityUUID,
		SkillType:    skillType,
		Points:       points,
	})
	if errors.Is(err, sm.ErrMaxPointsExceeded) {
		return c.Send("Превышено максимальное количество баллов, попробуй ещё раз.")
	} else if err != nil {
		return err
	}

	if err = c.Send(
		fmt.Sprintf(
			"Успешно начислено участнику @%s в навыки %q %d б.",
			username, skillType, points,
		),
	); err != nil {
		return err
	}

	return p.adminMenu(c, state)
}

func makeReplyButtonsApprove() [][]telebot.ReplyButton {
	buttons := []telebot.ReplyButton{
		{Text: approveButton},
		{Text: cancelButton},
	}
	return composeButtons(buttons, 2)
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

func makeReplyButtonsForSkillTypes(skills []string) [][]telebot.ReplyButton {
	buttons := make([]telebot.ReplyButton, len(skills))

	for i, skill := range skills {
		buttons[i] = telebot.ReplyButton{Text: skill}
	}

	return composeButtons(buttons, 2)
}

func makeReplyButtonsForRangeNumbers(start, end int) [][]telebot.ReplyButton {
	buttons := make([]telebot.ReplyButton, end-start+1)

	for i := start; i <= end; i++ {
		buttons[i-start] = telebot.ReplyButton{Text: strconv.Itoa(i)}
	}

	return composeButtons(buttons, 4)
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
