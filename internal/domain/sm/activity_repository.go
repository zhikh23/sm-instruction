package sm

import (
	"context"
	"errors"
)

var ErrActivityAlreadyExists = errors.New("activity already exists")
var ErrActivityNotFound = errors.New("activity not found")

type ActivitiesRepository interface {
	Save(ctx context.Context, activity *Activity) error
	Activity(ctx context.Context, activityName string) (*Activity, error)
	ActivityByAdmin(ctx context.Context, adminUsername string) (*Activity, error)
	AdditionalActivities(ctx context.Context) ([]*Activity, error)
	AvailableActivities(ctx context.Context) ([]*Activity, error)
	UpdateSlots(
		ctx context.Context,
		activityUUID string,
		updateFn func(innerCtx context.Context, activity *Activity) error,
	) error
}

type ActivitiesProvider interface {
	Activities(ctx context.Context) ([]*Activity, error)
}
