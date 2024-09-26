package telegram

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"gopkg.in/telebot.v3"

	"sm-instruction/internal/app/command"
	"sm-instruction/internal/app/query"
	"sm-instruction/internal/domain/sm"
)

const awardGroupNameKey = "awardGroupName"
const awardSkillKey = "awardSkill"

func (p *Port) awardSendEnterGroup(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	if err := s.SetState(ctx, awardHandleGroupNameState); err != nil {
		return err
	}

	return c.Send(buildMessage("\n",
		"Введи название учебной группы персонажа в формате:",
		"<code>СМ1-11Б</code>",
	),
		&telebot.ReplyMarkup{RemoveKeyboard: true}, telebot.ModeHTML,
	)
}

func (p *Port) awardHandleGroupName(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	groupName := c.Message().Text

	_, err := p.app.Queries.GetCharacter.Handle(ctx, query.GetCharacter{GroupName: groupName})
	if errors.Is(err, sm.ErrCharacterNotFound) {
		return p.awardSendCharacterNotFound(c, s)
	} else if err != nil {
		return err
	}

	if err = s.Update(ctx, awardGroupNameKey, groupName); err != nil {
		return err
	}

	return p.awardSendEnterSkill(c, s)
}

func (p *Port) awardSendEnterSkill(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	if err := s.SetState(ctx, awardHandleSkillState); err != nil {
		return err
	}

	activityName, err := extractActivityName(ctx, s)
	if err != nil {
		return err
	}

	act, err := p.app.Queries.GetActivity.Handle(ctx, query.GetActivity{ActivityName: activityName})
	if err != nil {
		return err
	}

	return c.Send(
		"Выбери один из доступных навыков для начисления баллов.",
		createMarkupWithButtonsFromStrings(act.Skills, 2),
	)
}

func (p *Port) awardSendCharacterNotFound(c telebot.Context, _ fsm.Context) error {
	return c.Send(buildMessage("\n",
		"Персонаж с такой группой не найден :(",
		"Проверь правильность вводимого формата:",
		"<code>СМ1-11Б</code>",
		"Попробуй ещё раз.",
	))
}

func (p *Port) awardHandleSkill(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	skillType := c.Message().Text

	activityName, err := extractActivityName(ctx, s)
	if err != nil {
		return err
	}

	act, err := p.app.Queries.GetActivity.Handle(ctx, query.GetActivity{ActivityName: activityName})
	if err != nil {
		return err
	}

	if !slices.Contains(act.Skills, skillType) {
		return p.awardSendInvalidSkill(c, s)
	}

	if err = s.SetState(ctx, awardHandlePointsState); err != nil {
		return err
	}

	if err = s.Update(ctx, awardSkillKey, skillType); err != nil {
		return err
	}

	return p.awardSendEnterPoints(c, s)
}

func (p *Port) awardSendEnterPoints(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	activityName, err := extractActivityName(ctx, s)
	if err != nil {
		return err
	}

	act, err := p.app.Queries.GetActivity.Handle(ctx, query.GetActivity{ActivityName: activityName})
	if err != nil {
		return err
	}

	options := make([]string, 0, act.MaxPoints+1)
	for i := 1; i <= act.MaxPoints; i++ {
		options = append(options, strconv.Itoa(i))
	}

	if err = s.SetState(ctx, awardHandlePointsState); err != nil {
		return err
	}

	return c.Send(
		"Выбери количество баллов для начисления.",
		createMarkupWithButtonsFromStrings(options, 4),
	)
}

func (p *Port) awardSendInvalidSkill(c telebot.Context, s fsm.Context) error {
	err := c.Send("Ой, кажется ты не можешь начислить баллы в этот навык. Попробуй ещё раз.")
	if err != nil {
		return err
	}
	return p.awardSendEnterSkill(c, s)
}

func (p *Port) awardHandlePoints(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	pointsStr := c.Message().Text
	points, err := strconv.Atoi(pointsStr)
	if err != nil {
		return err
	}

	activityName, err := extractActivityName(ctx, s)
	if err != nil {
		return err
	}

	groupName, err := awardExtractGroupName(ctx, s)
	if err != nil {
		return err
	}

	skill, err := awardExtractSkill(ctx, s)
	if err != nil {
		return err
	}

	err = p.app.Commands.AwardCharacter.Handle(ctx, command.AwardCharacter{
		GroupName:    groupName,
		ActivityName: activityName,
		SkillType:    skill,
		Points:       points,
	})
	if errors.Is(err, sm.ErrMaxPointsExceeded) {
		return p.awardSendInvalidPoints(c, s)
	} else if err != nil {
		return err
	}

	return p.awardSendSuccess(c, s, groupName, skill, points)
}

func (p *Port) awardSendSuccess(c telebot.Context, s fsm.Context, groupName string, skill string, points int) error {
	err := c.Send(
		fmt.Sprintf("Успешно начислены баллы %d в навыки %q группе %s", points, skill, groupName),
	)
	if err != nil {
		return err
	}
	return p.sendParticipantMenu(c, s)
}

func (p *Port) awardSendInvalidPoints(c telebot.Context, s fsm.Context) error {
	err := c.Send("Кажется это некорректное количество баллов.")
	if err != nil {
		return err
	}
	return p.awardSendEnterPoints(c, s)
}

func awardExtractGroupName(ctx context.Context, s fsm.Context) (string, error) {
	var groupName string
	if err := s.Data(ctx, awardGroupNameKey, &groupName); err != nil {
		return "", err
	}
	return groupName, nil
}

func awardExtractSkill(ctx context.Context, s fsm.Context) (string, error) {
	var skill string
	if err := s.Data(ctx, awardSkillKey, &skill); err != nil {
		return "", err
	}
	return skill, nil
}
