package telegram

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"gopkg.in/telebot.v3"

	"github.com/zhikh23/sm-instruction/internal/app/command"
	"github.com/zhikh23/sm-instruction/internal/app/query"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

const takeSlotActivityName = "takeSlotActivityName"
const takeSlotStartTime = "takeSlotStartTime"
const takeSlotApproveButton = "Да!"

func (p *Port) takeSlotSendChooseActivity(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	groupName, err := extractGroupName(ctx, s)
	if err != nil {
		return err
	}
	activities, err := p.app.Queries.AvailableActivities.Handle(ctx, query.AvailableActivities{
		GroupName: groupName,
	})
	if err != nil {
		return err
	}

	if len(activities) == 0 {
		return c.Send("🚫 Больше не осталось слотов для записи :(")
	}

	buttons := make([]string, len(activities))
	for i, activity := range activities {
		buttons[i] = activity.Name
	}

	msg := buildMessage(
		"<b>БРОНИРОВАНИЕ ТОЧКИ</b>\n",
		"Выбери доступную для записи точку.",
	)

	if err = s.SetState(ctx, takeSlotHandleActivityNameState); err != nil {
		return err
	}

	return c.Send(
		msg,
		createMarkupWithButtonsFromStrings(buttons, 2),
	)
}

func (p *Port) takeSlotHandleActivityName(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	activityName := c.Message().Text

	_, err := p.app.Queries.GetActivity.Handle(ctx, query.GetActivity{ActivityName: activityName})
	if errors.Is(err, sm.ErrActivityNotFound) {
		err = c.Send("🚫 Выбери одну из предложенных точек.")
		if err != nil {
			return err
		}
		return p.takeSlotSendChooseActivity(c, s)
	} else if err != nil {
		return err
	}

	if err = s.Update(ctx, takeSlotActivityName, activityName); err != nil {
		return err
	}

	return p.takeSlotSendSlots(c, s)
}

func (p *Port) takeSlotSendSlots(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	groupName, err := extractGroupName(ctx, s)
	if err != nil {
		return err
	}

	activityName, err := takeSkillExtractActivityName(ctx, s)
	if err != nil {
		return err
	}

	slots, err := p.app.Queries.AvailableSlots.Handle(ctx, query.AvailableSlots{
		GroupName:    groupName,
		ActivityName: activityName,
	})
	if err != nil {
		return err
	}

	if len(slots) == 0 {
		if err = c.Send("🚫 Нет свободных слотов для записи."); err != nil {
			return err
		}
		return p.sendParticipantMenu(c, s)
	}

	buttons := make([]string, len(slots))
	for i, slot := range slots {
		buttons[i] = slot.Start.Format(sm.TimeFormat)
	}

	if err = s.SetState(ctx, takeSlotHandleStartTimeState); err != nil {
		return err
	}

	return c.Send(
		"Выбери свободный промежуток времени.",
		createMarkupWithButtonsFromStrings(buttons, 4), telebot.ModeHTML,
	)
}

func (p *Port) takeSlotHandleStartTime(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	startS := c.Message().Text
	parsed, err := time.Parse(sm.TimeFormat, startS)
	if err != nil {
		return c.Send("🚫 Выбери корректное время начала точки.")
	}
	start := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), parsed.Hour(), parsed.Minute(), 0, 0, time.Local)

	if err = s.Update(ctx, takeSlotStartTime, start); err != nil {
		return err
	}

	return p.takeSlotSendActivity(c, s)
}

func (p *Port) takeSlotSendActivity(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	activityName, err := takeSkillExtractActivityName(ctx, s)
	if err != nil {
		return err
	}

	activity, err := p.app.Queries.GetActivity.Handle(ctx, query.GetActivity{
		ActivityName: activityName,
	})
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("<b>%s</b>\n", activity.FullName)

	if activity.Location != nil {
		msg = buildMessage("\n",
			msg,
			fmt.Sprintf("🔹 <i>Где?</i> %s", *activity.Location),
			"🔹 <i>Что это?</i>",
			"",
			*activity.Description,
			"",
			"❓ Подтверждаешь бронирование точки?",
		)
	}

	buttons := []string{takeSlotApproveButton, "Отменить"}

	if err = s.SetState(ctx, takeSlotHandleApproveState); err != nil {
		return err
	}

	return c.Send(
		msg, telebot.ModeHTML,
		createMarkupWithButtonsFromStrings(buttons, 2),
	)
}

func (p *Port) takeSlotHandleApprove(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	answer := c.Message().Text
	if answer != takeSlotApproveButton {
		return p.sendParticipantMenu(c, s)
	}

	groupName, err := extractGroupName(ctx, s)
	if err != nil {
		return err
	}

	activityName, err := takeSkillExtractActivityName(ctx, s)
	if err != nil {
		return err
	}

	start, err := takeSkillExtractStartTime(ctx, s)
	if err != nil {
		return err
	}

	err = p.app.Commands.TakeSlot.Handle(ctx, command.TakeSlot{
		GroupName:    groupName,
		ActivityName: activityName,
		Start:        start,
	})
	if errors.Is(err, sm.ErrSlotIsTooLate) {
		if err = c.Send("🚫 Ой, кажется ты пытаешься забронировать точку уже после окончания инструкции :("); err != nil {
			return err
		}
		return p.sendParticipantMenu(c, s)
	} else if err != nil {
		return err
	}

	err = c.Send(fmt.Sprintf("✅ Успешно забронирована точка %q на время %s", activityName, start.Format(sm.TimeFormat)))
	if err != nil {
		return err
	}

	return p.sendParticipantMenu(c, s)
}

func takeSkillExtractActivityName(ctx context.Context, s fsm.Context) (string, error) {
	var activityName string
	if err := s.Data(ctx, takeSlotActivityName, &activityName); err != nil {
		return "", fmt.Errorf("failed extract activity name: %w", err)
	}
	return activityName, nil
}

func takeSkillExtractStartTime(ctx context.Context, s fsm.Context) (time.Time, error) {
	var startTime time.Time
	if err := s.Data(ctx, takeSlotStartTime, &startTime); err != nil {
		return time.Time{}, fmt.Errorf("failed extract start time: %w", err)
	}
	return startTime, nil
}
