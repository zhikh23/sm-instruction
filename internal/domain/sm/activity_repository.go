package sm

import (
	"context"
	"errors"
)

var ErrActivityAlreadyExists = errors.New("activity already exists")
var ErrActivityNotFound = errors.New("activity not found")

type ActivitiesRepository interface {
	Save(ctx context.Context, activity *Activity) error
	Activity(ctx context.Context, activityUUID string) (*Activity, error)
	ActivityByName(ctx context.Context, name string) (*Activity, error)
	ActivityByAdmin(ctx context.Context, adminUsername string) (*Activity, error)
	ActivitiesWithLocations(ctx context.Context) ([]*Activity, error)
	Update(
		ctx context.Context,
		activityUUID string,
		updateFn func(innerCtx context.Context, activity *Activity) error,
	) error
}
