package app

import (
	"github.com/zhikh23/sm-instruction/internal/app/command"
	"github.com/zhikh23/sm-instruction/internal/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	StartInstruction command.StartInstructionHandler
	AwardCharacter   command.AwardCharacterHandler
	TakeSlot         command.TakeSlotHandler
}

type Queries struct {
	GetUser             query.GetUserHandler
	CharacterByUsername query.CharacterByUsernameHandler
	GetCharacter        query.GetCharacterHandler
	GetActivity         query.GetActivityHandler
	AdminActivity       query.AdminActivityHandler
	AvailableActivities query.AvailableActivitiesHandler
	AvailableSlots      query.AvailableSlotsHandler
}
