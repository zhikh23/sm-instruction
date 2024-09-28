package telegram

import (
	"context"
	"fmt"
	"slices"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/zhikh23/sm-instruction/internal/app/query"
	"gopkg.in/telebot.v3"
)

func (p *Port) sendParticipantRating(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	groupName, err := extractGroupName(ctx, s)
	if err != nil {
		return err
	}

	chars, err := p.app.Queries.Rating.Handle(ctx, query.Rating{})
	if err != nil {
		return err
	}

	current := slices.IndexFunc(chars, func(c query.Character) bool {
		return c.GroupName == groupName
	})

	msg := buildMessage("\n",
		"<b>СЕССИЯ</b>",
		"",
		"<i>ТОП-3</i>:",
		fmt.Sprintf("🥇 %d. %s - %0.2f", 1, chars[0].GroupName, chars[0].Rating),
		fmt.Sprintf("🥈 %d. %s - %0.2f", 2, chars[1].GroupName, chars[1].Rating),
		fmt.Sprintf("🥉 %d. %s - %0.2f", 3, chars[2].GroupName, chars[2].Rating),
	)
	if current > 2 {
		msg = buildMessage("\n",
			msg,
			"...",
			fmt.Sprintf("%d. %s - %0.2f", current+1, groupName, chars[current].Rating),
		)
	}

	if err = c.Send(msg, telebot.ModeHTML); err != nil {
		return err
	}

	return p.sendParticipantMenu(c, s)
}
