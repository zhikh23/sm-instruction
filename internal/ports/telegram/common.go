package telegram

import (
	"context"
	"fmt"
	"math"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"gopkg.in/telebot.v3"

	"github.com/zhikh23/sm-instruction/internal/app/query"
)

func (p *Port) sendIfError(c telebot.Context, _ fsm.Context) error {
	return c.Send(
		"Что-то пошло не так...\n" +
			"Пожалуйста, обратись к организатором о случившийся проблеме и " +
			"не забудь сообщить свой ник: @" + c.Chat().Username + ".",
	)
}

func (p *Port) sendIfErrorDebug(c telebot.Context, _ fsm.Context, err error) error {
	return c.Send(err.Error())
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

func createMarkupWithButtonsFromStrings(strings []string, width int) *telebot.ReplyMarkup {
	buttons := make([]telebot.ReplyButton, len(strings))
	for i, text := range strings {
		buttons[i] = telebot.ReplyButton{Text: text}
	}
	return &telebot.ReplyMarkup{
		ReplyKeyboard: composeButtons(buttons, width),
	}
}

const (
	participantMenuProfileButton    = "Профиль"
	participantMenuTimetableButton  = "Расписание"
	participantMenuTakeSlotButton   = "Забронировать точку"
	participantMenuGradesButton     = "Успеваемость"
	participantMenuRatingButton     = "Сессия"
	participantMenuAdditionalButton = "Дополнительные задания"

	adminMenuAwardCharacterButton = "Начислить баллы"
	adminMenuTimetableButton      = "Расписание"
)

func (p *Port) sendParticipantMenu(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	if err := s.SetState(ctx, participantMenuHandle); err != nil {
		return err
	}

	return c.Send(
		"Выбери действие.",
		createMarkupWithButtonsFromStrings([]string{
			participantMenuProfileButton,
			participantMenuTimetableButton,
			participantMenuTakeSlotButton,
			participantMenuGradesButton,
			participantMenuRatingButton,
			participantMenuAdditionalButton,
		}, 2),
	)
}

func (p *Port) sendAdminMenu(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	if err := s.SetState(ctx, adminMenuHandle); err != nil {
		return err
	}

	act, err := p.app.Queries.AdminActivity.Handle(ctx, query.AdminActivity{Username: c.Chat().Username})
	if err != nil {
		return err
	}

	if err = s.Update(ctx, activityNameKey, act.Name); err != nil {
		return err
	}

	return c.Send(
		"Панель управления администратора.",
		createMarkupWithButtonsFromStrings([]string{
			adminMenuAwardCharacterButton,
			adminMenuTimetableButton,
		}, 2),
	)
}

func extractGroupName(ctx context.Context, s fsm.Context) (string, error) {
	var groupName string
	if err := s.Data(ctx, groupNameKey, &groupName); err != nil {
		return "", fmt.Errorf("failed extract group name: %w", err)
	}
	return groupName, nil
}

func extractActivityName(ctx context.Context, s fsm.Context) (string, error) {
	var activityName string
	if err := s.Data(ctx, activityNameKey, &activityName); err != nil {
		return "", fmt.Errorf("failed extract activity name: %w", err)
	}
	return activityName, nil
}
