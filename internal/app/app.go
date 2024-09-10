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
	CancelBooking command.CancelBookingHandler
}

type Queries struct {
	GetCharacter query.GetCharacterHandler

	GetLocation           query.GetLocationHandler
	GetAllLocations       query.GetAllLocationsHandler
	GetLocationByName     query.GetLocationByNameHandler
	GetAvailableIntervals query.GetAvailableIntervalsHandler
}
