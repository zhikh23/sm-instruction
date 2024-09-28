package telegram

import (
	"context"
	"errors"
	"fmt"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"gopkg.in/telebot.v3"

	"github.com/zhikh23/sm-instruction/internal/app/query"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

func (p *Port) sendParticipantAdditionalActivities(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	groupName, err := extractGroupName(ctx, s)
	if err != nil {
		return err
	}
	activities, err := p.app.Queries.AdditionalActivities.Handle(ctx, query.AdditionalActivities{
		GroupName: groupName,
	})
	if err != nil {
		return err
	}

	msg := buildMessage("\n",
		"<b>ДОПОЛНИТЕЛЬНЫЕ ЗАДАНИЯ</b>",
		"",
		"Здесь можно узнать подробнее о дополнительных активностях.",
	)

	buttons := make([]string, len(activities))
	for i, activity := range activities {
		buttons[i] = activity.Name
	}

	if err = s.SetState(ctx, additionalHandleActivityNameState); err != nil {
		return err
	}

	return c.Send(
		msg, telebot.ModeHTML,
		createMarkupWithButtonsFromStrings(buttons, 2),
	)
}

func (p *Port) additionalHandleActivityName(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	activityName := c.Message().Text

	activity, err := p.app.Queries.GetActivity.Handle(ctx, query.GetActivity{ActivityName: activityName})
	if errors.Is(err, sm.ErrActivityNotFound) {
		err = c.Send("🚫 Выбери одно из предложенных дополнительных заданий.")
		if err != nil {
			return err
		}
		return p.sendParticipantAdditionalActivities(c, s)
	} else if err != nil {
		return err
	}

	msg := buildMessage("\n",
		fmt.Sprintf("<b>%s</b>", activity.Name),
		"",
		fmt.Sprintf("🔹 %s", *activity.Description),
	)

	if err = c.Send(msg, telebot.ModeHTML); err != nil {
		return err
	}

	return p.sendParticipantMenu(c, s)
}
