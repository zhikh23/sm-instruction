package telegram

import (
	"context"
	"fmt"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"gopkg.in/telebot.v3"

	"github.com/zhikh23/sm-instruction/internal/app/query"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
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

	msg := buildMessage("\n",
		"<b>РАСПИСАНИЕ</b>\n",
		fmt.Sprintf("Для группы %s:", groupName),
		"",
	)
	for _, slot := range char.Slots {
		act, err := p.app.Queries.GetActivity.Handle(ctx, query.GetActivity{ActivityName: *slot.Whom})
		if err != nil {
			return err
		}

		if slot.Whom == nil {
			continue
		}
		msg = buildMessage("\n",
			msg,
			fmt.Sprintf(
				"<code>%s-%s</code> | %s | %s",
				slot.Start.Format(sm.TimeFormat), slot.End.Format(sm.TimeFormat), *slot.Whom, *act.Location,
			),
		)
	}

	err = c.Send(msg, &telebot.ReplyMarkup{RemoveKeyboard: true}, telebot.ModeHTML)
	if err != nil {
		return err
	}

	return p.sendParticipantMenu(c, s)
}
