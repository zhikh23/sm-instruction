package telegram

import (
	"context"
	"errors"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"gopkg.in/telebot.v3"

	"github.com/zhikh23/sm-instruction/internal/app/command"
	"github.com/zhikh23/sm-instruction/internal/app/query"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

func (p *Port) StartHandleCommand(c telebot.Context, s fsm.Context) error {
	ctx := context.Background()

	user, err := p.app.Queries.GetUser.Handle(ctx, query.GetUser{Username: c.Chat().Username})
	if errors.Is(err, sm.ErrUserNotFound) {
		return p.sendUserNotFound(c, s)
	} else if err != nil {
		return err
	}

	if user.Role == "administrator" {
		act, err := p.app.Queries.AdminActivity.Handle(ctx, query.AdminActivity{Username: c.Chat().Username})
		if err != nil {
			return err
		}

		if err = s.Update(ctx, activityNameKey, act.Name); err != nil {
			return err
		}

		return p.sendAdminMenu(c, s)
	}

	char, err := p.app.Queries.CharacterByUsername.Handle(ctx, query.CharacterByUsername{
		Username: c.Chat().Username,
	})
	if err != nil {
		return err
	}

	if err = s.Update(ctx, groupNameKey, char.GroupName); err != nil {
		return err
	}

	err = p.app.Commands.StartInstruction.Handle(ctx, command.StartInstruction{
		GroupName: char.GroupName,
	})
	if err != nil {
		return err
	}

	err = p.sendParticipantStartMessage(c, s)
	if err != nil {
		return err
	}

	return p.sendParticipantMenu(c, s)
}

func (p *Port) sendUserNotFound(c telebot.Context, _ fsm.Context) error {
	return c.Send(
		"🚫 Я тебя ещё не знаю :(\n" +
			"Пожалуйста, обратись к организаторам «СМ. Инструкции по выживанию».",
	)
}

func (p *Port) sendParticipantStartMessage(c telebot.Context, _ fsm.Context) error {
	if _, err := studentSticker.Send(c.Bot(), c.Recipient(), nil); err != nil {
		return err
	}
	return c.Send(
		"Привет, участник «СМ. Инструкция по выживанию»!",
	)
}
