package telegram

import (
	"context"
	"fmt"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"gopkg.in/telebot.v3"

	"github.com/zhikh23/sm-instruction/internal/app/query"
)

func (p *Port) sendCharacterTimetable(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	groupName, err := extractGroupName(ctx, s)
	if err != nil {
		return err
	}

	char, err := p.app.Queries.GetCharacter.Handle(ctx, query.GetCharacter{GroupName: groupName})
	if err != nil {
		return err
	}

	msg := "<b>РАСПИСАНИЕ</b>\n\n"
	for _, slot := range char.Slots {
		var text string
		if slot.Whom == nil {
			text = "-"
		} else {
			text = *slot.Whom
		}
		msg += fmt.Sprintf(
			"<code>%s-%s</code> %s\n",
			slot.Start.Format(sm.TimeFormat), slot.End.Format(sm.TimeFormat), text,
		)
	}

	err = c.Send(msg, &telebot.ReplyMarkup{RemoveKeyboard: true}, telebot.ModeHTML)
	if err != nil {
		return err
	}

	return p.sendParticipantMenu(c, s)
}
