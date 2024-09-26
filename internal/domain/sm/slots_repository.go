package sm

import (
	"context"
	"errors"
	"time"
)

var ErrSlotAlreadyExists = errors.New("slot already exists")

type BookingsRepository interface {
	Save(ctx context.Context, slot *Slot) error
	Update(
		ctx context.Context,
		activityName string,
		startTime time.Time,
		updateFn func(innerCtx context.Context) error,
	) error

	BookedSlots(ctx context.Context, groupName string) ([]*Slot, error)
	ActivitySlots(ctx context.Context, activityName string) ([]*Slot, error)
	AvailableSlots(ctx context.Context, activityName string) ([]*Slot, error)
}
