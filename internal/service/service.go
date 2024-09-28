package service

import (
	"errors"
	"log/slog"

	"github.com/zhikh23/sm-instruction/internal/adapters"
	"github.com/zhikh23/sm-instruction/internal/app"
	"github.com/zhikh23/sm-instruction/internal/app/command"
	"github.com/zhikh23/sm-instruction/internal/app/query"
	"github.com/zhikh23/sm-instruction/internal/common/decorator"
	"github.com/zhikh23/sm-instruction/internal/common/logs"
	"github.com/zhikh23/sm-instruction/internal/common/metrics"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
	"github.com/zhikh23/sm-instruction/internal/service/mocks"
)

func NewApplication() (*app.Application, func() error) {
	log := logs.DefaultLogger()
	metricsClient := metrics.NoOp{}

	users, closeUsers := adapters.NewPGUsersRepository()
	chars, closeChars := adapters.NewPGCharactersRepository()
	activities, closeActivities := adapters.NewPGActivitiesRepository()

	return newApplication(log, metricsClient, users, chars, activities), func() error {
		var err error
		err = errors.Join(err, closeUsers())
		err = errors.Join(err, closeChars())
		err = errors.Join(err, closeActivities())
		return err
	}
}

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
			Rating:              query.NewRatingHandler(chars, log, metricsClient),
			GetActivity:         query.NewGetActivityHandler(activities, log, metricsClient),
			AdminActivity:       query.NewAdminActivtyHandler(activities, log, metricsClient),
			AvailableActivities: query.NewAvailableActivitiesHandler(chars, activities, log, metricsClient),
			AvailableSlots:      query.NewAvailableSlotsHandler(chars, activities, log, metricsClient),
		},
	}
}
