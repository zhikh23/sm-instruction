package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"gopkg.in/telebot.v3"

	"github.com/zhikh23/sm-instruction/internal/app/query"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
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
		"<b>ПРОФИЛЬ</b>",
		"",
		fmt.Sprintf("Учебная группа: <code>%s</code>", char.GroupName),
	)
	if time.Now().Before(*char.End) {
		remains := char.End.Sub(*char.Start)
		msg = buildMessage("\n",
			msg,
			"",
			fmt.Sprintf("Начало Инструкции: %s", char.Start.Format(sm.TimeFormat)),
			fmt.Sprintf("Конец Инструкции: %s", char.End.Format(sm.TimeFormat)),
			fmt.Sprintf(
				"❕ Осталось до конца Инструкции: <b>%d:%02d</b>\n",
				int(remains.Hours()), int(remains.Minutes())%60,
			),
		)
	} else {
		msg = buildMessage("\n",
			msg,
			fmt.Sprintf("Инструкция окончена в <b>%s</b>.", char.Start.Format(sm.TimeFormat)),
		)
	}

	msg = buildMessage("\n",
		msg,
		"<b>Навыки:</b>",
		fmt.Sprintf("🛠 <i>Инженерные - %d</i>", char.Skills[sm.Engineering.String()]),
		fmt.Sprintf("🔭 <i>Исследовательские - %d</i>", char.Skills[sm.Researching.String()]),
		fmt.Sprintf("🤝 <i>Социальные - %d</i>", char.Skills[sm.Social.String()]),
		fmt.Sprintf("⚽️ <i>Спортивные - %d</i>", char.Skills[sm.Sportive.String()]),
		fmt.Sprintf("🔮 <i>Творческие - %d</i>", char.Skills[sm.Creative.String()]),
		"",
		fmt.Sprintf("🏅 Рейтинг: <b>%0.1f</b>", char.Rating),
	)

	if err = c.Send(msg, telebot.ModeHTML); err != nil {
		return err
	}

	return p.sendParticipantMenu(c, s)
}

func buildMessage(sep string, lines ...string) string {
	return strings.Join(lines, sep)
}
