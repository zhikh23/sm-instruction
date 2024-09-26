package telegram

import (
	"context"
	"fmt"
	"sm-instruction/internal/domain/sm"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"gopkg.in/telebot.v3"

	"sm-instruction/internal/app/query"
)

func (p *Port) sendAdminTimetable(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	activityName, err := extractActivityName(ctx, s)
	if err != nil {
		return err
	}

	char, err := p.app.Queries.GetActivity.Handle(ctx, query.GetActivity{ActivityName: activityName})
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Расписание на точку %q «СМ. Инструкция по выживанию»:\n\n", activityName)
	for _, slot := range char.Slots {
		var text string
		if slot.Whom == nil {
			text = "Свободно"
		} else {
			text = *slot.Whom
		}
		msg += fmt.Sprintf(
			"<code>%s-%s</code> %q\n",
			slot.Start.Format(sm.TimeFormat), slot.End.Format(sm.TimeFormat), text,
		)
	}

	err = c.Send(msg, &telebot.ReplyMarkup{RemoveKeyboard: true}, telebot.ModeHTML)
	if err != nil {
		return err
	}

	return p.sendAdminMenu(c, s)
}
