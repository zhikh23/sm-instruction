package sm

import (
	"context"
	"errors"
)

var ErrLocationAlreadyExists = errors.New("location already exists")
var ErrLocationNotFound = errors.New("location not found")

type LocationsRepository interface {
	Save(ctx context.Context, l *Location) error
	Location(ctx context.Context, uuid string) (*Location, error)
	LocationByName(ctx context.Context, name string) (*Location, error)
	Locations(ctx context.Context) ([]*Location, error)
	Update(
		ctx context.Context,
		locationUUID string,
		updateFn func(innerCtx context.Context, loc *Location) error,
	) error
}
