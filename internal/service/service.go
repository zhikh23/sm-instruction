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
			AwardCharacter:   command.NewAwardCharacterHandler(chars, activities, log, metricsClient),
			TakeSlot:         command.NewTakeSlotHandler(chars, activities, log, metricsClient),
		},
		Queries: app.Queries{
			GetUser:             query.NewGetUserHandler(users, log, metricsClient),
			CharacterByUsername: query.NewCharacterByUsernameHandler(chars, log, metricsClient),
			GetCharacter:        query.NewGetCharacterHandler(chars, log, metricsClient),
			GetActivity:         query.NewGetActivityHandler(activities, log, metricsClient),
			AdminActivity:       query.NewAdminActivtyHandler(activities, log, metricsClient),
			AvailableActivities: query.NewAvailableActivitiesHandler(activities, log, metricsClient),
			AvailableSlots:      query.NewAvailableSlotsHandler(chars, activities, log, metricsClient),
		},
	}
}
