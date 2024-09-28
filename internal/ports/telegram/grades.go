package telegram

import (
	"context"
	"fmt"
	"strconv"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"gopkg.in/telebot.v3"

	"github.com/zhikh23/sm-instruction/internal/app/query"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

func (p *Port) sendParticipantsGrades(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	groupName, err := extractGroupName(ctx, s)
	if err != nil {
		return err
	}

	char, err := p.app.Queries.GetCharacter.Handle(ctx, query.GetCharacter{GroupName: groupName})
	if err != nil {
		return err
	}

	msg := "<b>УСПЕВАЕМОСТЬ</b>\n\n"
	for _, grade := range char.Grades {
		msg += buildMessage(" ",
			grade.Time.Format(sm.TimeFormat),
			"-",
			fmt.Sprintf("%q", grade.ActivityName),
			"-",
			grade.SkillType,
			"-",
			"<i>"+strconv.Itoa(grade.Points),
			"б</i>\n",
		)
	}

	if err = c.Send(msg, telebot.ModeHTML); err != nil {
		return err
	}

	return p.sendParticipantMenu(c, s)
}
