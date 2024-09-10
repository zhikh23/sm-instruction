package service

import (
	"log/slog"

	"sm-instruction/internal/app"
	"sm-instruction/internal/app/command"
	"sm-instruction/internal/app/query"
	"sm-instruction/internal/common/decorator"
	"sm-instruction/internal/common/logs"
	"sm-instruction/internal/common/metrics"
	"sm-instruction/internal/domain/sm"
	"sm-instruction/internal/service/mocks"
)

func NewMockedApplication() *app.Application {
	log := logs.DefaultLogger()
	metricsClient := metrics.NoOp{}

	chars := mocks.NewMockCharactersRepository()
	locs := mocks.NewMockLocationsRepository()

	return newApplication(log, metricsClient, chars, locs)
}

func newApplication(
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
	chars sm.CharactersRepository,
	locs sm.LocationsRepository,
) *app.Application {
	return &app.Application{
		Commands: app.Commands{
			StartInstruction: command.NewStartInstructionHandler(chars, log, metricsClient),

			BookLocation:  command.NewBookLocationHandler(chars, locs, log, metricsClient),
			CancelBooking: command.NewCancelBookingHandler(chars, locs, log, metricsClient),
		},
		Queries: app.Queries{
			GetCharacter: query.NewGetCharacterHandler(chars, log, metricsClient),

			GetLocation:           query.NewGetLocationHandler(locs, log, metricsClient),
			GetAllLocations:       query.NewGetAllLocationsHandler(locs, log, metricsClient),
			GetLocationByName:     query.NewGetLocationByNameHandler(locs, log, metricsClient),
			GetAvailableIntervals: query.NewGetAvailableIntervalsHandler(chars, locs, log, metricsClient),
		},
	}
}
