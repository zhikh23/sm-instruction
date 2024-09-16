package app

import (
	"sm-instruction/internal/app/command"
	"sm-instruction/internal/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	StartInstruction command.StartInstructionHandler

	BookLocation  command.BookLocationHandler
	RemoveBooking command.RemoveBookingHandler

	AwardCharacter command.AwardCharacterHandler
}

type Queries struct {
	UserIsAdministrator query.UserIsAdministratorHandler

	CharacterIsStarted  query.CharacterIsStartedHandler
	GetCharacter        query.GetCharacterHandler
	GetCharacterByGroup query.GetCharacterByGroupHandler

	GetActivity        query.GetActivityHandler
	GetActivityByAdmin query.GetActivityByAdminHandler

	GetLocation           query.GetLocationHandler
	GetAllLocations       query.GetAllLocationsHandler
	GetLocationByName     query.GetLocationByNameHandler
	GetAvailableIntervals query.GetAvailableIntervalsHandler
}
