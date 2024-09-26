package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"gopkg.in/telebot.v3"

	"sm-instruction/internal/app/query"
	"sm-instruction/internal/domain/sm"
)

func (p *Port) sendProfile(c telebot.Context, s fsm.Context) error {
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
		fmt.Sprintf("Твой никнейм: @%s", char.Username),
		"Твои навыки:",
	)
	for _, skill := range sm.AllSkills {
		msg = buildMessage("\n", msg,
			fmt.Sprintf("<i>%s</i> - %d", skill.String(), char.Skills[skill.String()]),
		)
	}
	msg = buildMessage("\n", msg,
		fmt.Sprintf("Твой рейтинг: <b>%0.1f</b>", char.Rating),
	)

	if err = c.Send(msg, telebot.ModeHTML); err != nil {
		return err
	}

	return p.sendParticipantMenu(c, s)
}

func buildMessage(sep string, lines ...string) string {
	return strings.Join(lines, sep)
}
