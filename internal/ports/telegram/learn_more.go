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

const learnMoreActivityNameKey = "learnMoreActivityName"

func (p *Port) learnMoreSendActivities(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	activities, err := p.app.Queries.Activities.Handle(ctx, query.Activities{})
	if err != nil {
		return err
	}

	msg := buildMessage("\n",
		"<b>МАТЕРИАЛЫ</b>",
		"",
		"Вся-вся информация в одной вкладке!",
	)

	buttons := make([]string, len(activities))
	for i, activity := range activities {
		buttons[i] = activity.Name
	}

	if err = s.SetState(ctx, learnMoreHandleActivityNameState); err != nil {
		return err
	}

	return c.Send(msg, telebot.ModeHTML, createMarkupWithButtonsFromStrings(buttons, 4))
}

func (p *Port) learnMoreHandleActivityName(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	activityName := c.Message().Text

	if err := s.Update(ctx, learnMoreActivityNameKey, activityName); err != nil {
		return err
	}

	return p.learnMoreSendActivity(c, s)
}

func (p *Port) learnMoreSendActivity(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	activityName, err := learnMoreExtractActivityName(ctx, s)
	if err != nil {
		return err
	}

	activity, err := p.app.Queries.GetActivity.Handle(ctx, query.GetActivity{
		ActivityName: activityName,
	})
	if errors.Is(err, sm.ErrActivityNotFound) {
		err = c.Send("🚫 Выбери одну из предложенных точек.")
		if err != nil {
			return err
		}
		return p.learnMoreSendActivities(c, s)
	} else if err != nil {
		return err
	}

	msg := fmt.Sprintf("<b>%s</b>\n", activityName)

	if activity.Location != nil {
		msg = buildMessage("\n",
			msg,
			fmt.Sprintf("🔹 <i>Где?</i> %s", *activity.Location),
		)
	}

	if activity.Description != nil {
		msg = buildMessage("\n",
			msg,
			"🔹 <i>Что это?</i>",
			"",
			*activity.Description,
		)
	}

	if err = c.Send(msg, telebot.ModeHTML); err != nil {
		return err
	}

	return p.sendParticipantMenu(c, s)
}

func learnMoreExtractActivityName(ctx context.Context, s fsm.Context) (string, error) {
	var activityName string
	if err := s.Data(ctx, learnMoreActivityNameKey, &activityName); err != nil {
		return "", fmt.Errorf("failed extract activity name: %w", err)
	}
	return activityName, nil
}
