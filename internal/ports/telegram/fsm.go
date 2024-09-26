package telegram

import (
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/fsmopt"
	"gopkg.in/telebot.v3"
)

const groupNameKey = "groupName"
const activityNameKey = "groupActivityName"

const (
	participantMenuHandle = fsm.State("participantMenuHandle")
	adminMenuHandle       = fsm.State("adminMenuHandle")

	awardHandleGroupNameState = fsm.State("awardHandleGroupNameState")
	awardHandleSkillState     = fsm.State("awardHandleSkillState")
	awardHandlePointsState    = fsm.State("awardHandlePointsState")

	takeSlotHandleActivityNameState = fsm.State("takeSlotHandleActivityNameState")
	takeSlotHandleStartTimeState    = fsm.State("takeSlotHandleStartTimeState")
)

func (p *Port) RegisterFSMManager(m *fsm.Manager, dp fsm.Dispatcher) {
	dp.Dispatch(m.New(
		fsmopt.OnStates(fsm.AnyState),
		fsmopt.On("/start"),
		fsmopt.Do(p.StartHandleCommand),
	))

	dp.Dispatch(m.New(
		fsmopt.OnStates(fsm.AnyState),
		fsmopt.On("/admin"),
		fsmopt.Do(p.sendAdminMenu),
	))

	dp.Dispatch(m.New(
		fsmopt.OnStates(participantMenuHandle),
		fsmopt.On(participantMenuProfileButton),
		fsmopt.Do(p.sendProfile),
	))

	dp.Dispatch(m.New(
		fsmopt.OnStates(participantMenuHandle),
		fsmopt.On(participantMenuTimetableButton),
		fsmopt.Do(p.sendCharacterTimetable),
	))

	dp.Dispatch(m.New(
		fsmopt.OnStates(participantMenuHandle),
		fsmopt.On(participantMenuTakeSlotButton),
		fsmopt.Do(p.takeSlotSendChooseActivity),
	))

	dp.Dispatch(m.New(
		fsmopt.OnStates(adminMenuHandle),
		fsmopt.On(adminMenuTimetableButton),
		fsmopt.Do(p.sendAdminTimetable),
	))

	dp.Dispatch(m.New(
		fsmopt.OnStates(adminMenuHandle),
		fsmopt.On(adminMenuAwardCharacterButton),
		fsmopt.Do(p.awardSendEnterGroup),
	))

	dp.Dispatch(m.New(
		fsmopt.OnStates(awardHandleGroupNameState),
		fsmopt.On(telebot.OnText),
		fsmopt.Do(p.awardHandleGroupName),
	))

	dp.Dispatch(m.New(
		fsmopt.OnStates(awardHandleSkillState),
		fsmopt.On(telebot.OnText),
		fsmopt.Do(p.awardHandleSkill),
	))

	dp.Dispatch(m.New(
		fsmopt.OnStates(awardHandlePointsState),
		fsmopt.On(telebot.OnText),
		fsmopt.Do(p.awardHandlePoints),
	))

	dp.Dispatch(m.New(
		fsmopt.OnStates(takeSlotHandleActivityNameState),
		fsmopt.On(telebot.OnText),
		fsmopt.Do(p.takeSlotHandleActivityName),
	))

	dp.Dispatch(m.New(
		fsmopt.OnStates(takeSlotHandleStartTimeState),
		fsmopt.On(telebot.OnText),
		fsmopt.Do(p.takeSlotHandleStartTime),
	))
}
