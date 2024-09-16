package sm

import (
	"errors"
	"time"

	"sm-instruction/internal/common/commonerrs"
)

type Location struct {
	Description string
	Where       string
	BookedTimes []BookedTime
}

func NewLocation(
	description string,
	where string,
) (*Location, error) {
	if description == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty Location description")
	}

	if where == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty Location where")
	}

	return &Location{
		Description: description,
		Where:       where,
		BookedTimes: make([]BookedTime, 0),
	}, nil
}

func MustNewLocation(
	description string,
	where string,
) *Location {
	l, err := NewLocation(description, where)
	if err != nil {
		panic(err)
	}
	return l
}

func UnmarshallLocationFromDB(
	activityUUID string,
	name string,
	description string,
	where string,
	bookedTimes []BookedTime,
) (*Location, error) {
	if activityUUID == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty Location uuid")
	}

	if name == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty Location name")
	}

	if description == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty Location description")
	}

	if where == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty Location where")
	}

	if bookedTimes == nil {
		bookedTimes = make([]BookedTime, 0)
	}

	for _, interval := range bookedTimes {
		if interval.IsZero() {
			return nil, commonerrs.NewInvalidInputError("expected not empty booking interval")
		}
	}

	return &Location{
		Description: description,
		Where:       where,
		BookedTimes: bookedTimes,
	}, nil
}

func (l *Location) IsBooked(t time.Time) bool {
	for _, bt := range l.BookedTimes {
		if bt.Start.Equal(t) {
			return true
		}
	}
	return false
}

var ErrLocationAlreadyBooked = errors.New("interval already booked")

func (l *Location) CanBook(t time.Time) error {
	if l.IsBooked(t) {
		return ErrLocationAlreadyBooked
	}
	return nil
}

func (l *Location) AddBooking(bookedTime BookedTime) error {
	if err := l.CanBook(bookedTime.Start); err != nil {
		return err
	}

	l.BookedTimes = append(l.BookedTimes, bookedTime)

	return nil
}

var ErrLocationHasNotBooked = errors.New("interval is not booked")

func (l *Location) RemoveBooking(bookedTime BookedTime) error {
	for i, bt := range l.BookedTimes {
		if bt == bookedTime {
			l.BookedTimes = append(l.BookedTimes[:i], l.BookedTimes[i+1:]...)
			return nil
		}
	}
	return ErrLocationHasNotBooked
}

func (l *Location) AvailableTimes(to time.Time) []time.Time {
	times := make([]time.Time, 0)

	it := time.Now().Round(BookInterval)
	if it.Before(time.Now()) {
		it = it.Add(BookInterval)
	}

	for ; it.Before(to); it = it.Add(BookInterval) {
		if !l.IsBooked(it) {
			times = append(times, it)
		}
	}

	return times
}
