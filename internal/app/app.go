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
	BookLocation     command.BookLocationHandler
}

type Queries struct {
	GetCharacter          query.GetCharacterHandler
	GetAllLocations       query.GetAllLocationsHandler
	GetLocationByName     query.GetLocationByNameHandler
	GetAvailableIntervals query.GetAvailableIntervalsHandler
}
