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

	users := mocks.NewMockUsersRepository()
	chars := mocks.NewMockCharactersRepository()
	activities := mocks.NewMockActivitiesRepository()

	return newApplication(log, metricsClient, users, chars, activities)
}

func newApplication(
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
	users sm.UsersRepository,
	chars sm.CharactersRepository,
	activities sm.ActivitiesRepository,
) *app.Application {
	return &app.Application{
		Commands: app.Commands{
			StartInstruction: command.NewStartInstructionHandler(users, chars, log, metricsClient),

			BookLocation:  command.NewBookLocationHandler(chars, activities, log, metricsClient),
			CancelBooking: command.NewRemoveBookingHandler(chars, activities, log, metricsClient),
		},
		Queries: app.Queries{
			UserIsAdministrator: query.NewUserIsAdministratorHandler(users, log, metricsClient),

			CharacterIsStarted: query.NewCharacterIsStarted(chars, log, metricsClient),
			GetCharacter:       query.NewGetCharacterHandler(chars, log, metricsClient),

			GetLocation:           query.NewGetLocationHandler(activities, log, metricsClient),
			GetAllLocations:       query.NewGetAllLocationsHandler(activities, log, metricsClient),
			GetLocationByName:     query.NewGetLocationByNameHandler(activities, log, metricsClient),
			GetActivityByAdmin:    query.NewGetLocationByAdminHandler(activities, log, metricsClient),
			GetAvailableIntervals: query.NewGetAvailableIntervalsHandler(chars, activities, log, metricsClient),
		},
	}
}
