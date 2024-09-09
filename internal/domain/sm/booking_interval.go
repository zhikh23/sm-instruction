package sm

import (
	"errors"
	"time"

	"sm-instruction/internal/common/commonerrs"
)

type BookingInterval struct {
	From     time.Time
	To       time.Time
	BookedBy string
}

func UnmarshalBookingIntervalFromDB(
	from time.Time,
	to time.Time,
	bookedBy string,
) (BookingInterval, error) {
	if from.IsZero() {
		return BookingInterval{}, commonerrs.NewInvalidInputError("expected not empty 'from' timestamp")
	}

	if to.IsZero() {
		return BookingInterval{}, commonerrs.NewInvalidInputError("expected not empty 'to' timestamp")
	}

	if bookedBy == "" {
		return BookingInterval{}, commonerrs.NewInvalidInputError("expected not empty booked by ID")
	}

	return BookingInterval{
		From:     from.Local(),
		To:       to.Local(),
		BookedBy: bookedBy,
	}, nil
}

func (i BookingInterval) IsZero() bool {
	return i == BookingInterval{}
}

func (i BookingInterval) IsIntersects(o BookingInterval) bool {
	return i.From.Before(o.To) && i.To.After(o.From)
}

type BookingIntervalFactoryConfig struct {
	IntervalDuration             time.Duration
	MinimalDurationBeforeBooking time.Duration
}

func (c BookingIntervalFactoryConfig) Validate() error {
	var err error

	if c.IntervalDuration <= 0 {
		err = errors.Join(err, commonerrs.NewInvalidInputError(
			"booking interval should be positive duration",
		))
	}

	if c.MinimalDurationBeforeBooking <= 0 {
		err = errors.Join(err, commonerrs.NewInvalidInputError(
			"minimal booking interval should be positive duration",
		))
	}

	return err
}

type BookingIntervalFactory struct {
	cfg BookingIntervalFactoryConfig
}

func NewBookingIntervalFactory(cfg BookingIntervalFactoryConfig) (BookingIntervalFactory, error) {
	if err := cfg.Validate(); err != nil {
		return BookingIntervalFactory{}, err
	}

	return BookingIntervalFactory{cfg: cfg}, nil
}

func MustNewBookingIntervalFactory(cfg BookingIntervalFactoryConfig) BookingIntervalFactory {
	f, err := NewBookingIntervalFactory(cfg)
	if err != nil {
		panic(err)
	}
	return f
}

var ErrNowBookingIsTooLate = errors.New("too late to book")

func (f BookingIntervalFactory) NewBookingInterval(from time.Time, username string) (BookingInterval, error) {
	if !from.Round(f.cfg.IntervalDuration).Equal(from) {
		return BookingInterval{}, commonerrs.NewInvalidInputErrorf(
			"invalid booking interval: should be multyply of %s", f.cfg.IntervalDuration.String(),
		)
	}

	if from.Sub(time.Now()) < f.cfg.MinimalDurationBeforeBooking {
		return BookingInterval{}, ErrNowBookingIsTooLate
	}

	return BookingInterval{
		From:     from.Local(),
		To:       from.Add(f.cfg.IntervalDuration).Local(),
		BookedBy: username,
	}, nil
}

func (f BookingIntervalFactory) MustNewBookingInterval(from time.Time, username string) BookingInterval {
	i, err := f.NewBookingInterval(from, username)
	if err != nil {
		panic(err)
	}
	return i
}

func (f BookingIntervalFactory) AvailableIntervals(char *Character, loc *Location) ([]BookingInterval, error) {
	finish, err := char.FinishTime()
	if err != nil {
		return nil, err
	}

	intervals, err := f.availableIntervals(finish, char.Username())
	if err != nil {
		return nil, err
	}

	availableIntervals := make([]BookingInterval, 0, len(intervals))
	for _, interval := range intervals {
		if loc.IsBooked(interval) {
			continue
		}
		if err = char.CanBook(loc, interval); err != nil {
			continue
		}
		availableIntervals = append(availableIntervals, interval)
	}

	return availableIntervals, nil
}

func (f BookingIntervalFactory) availableIntervals(finish time.Time, username string) ([]BookingInterval, error) {
	start := time.Now().Round(f.cfg.IntervalDuration)
	if start.Before(time.Now()) {
		start = start.Add(f.cfg.IntervalDuration)
	}

	availableIntervals := make([]BookingInterval, 0)
	for current := start; current.Before(finish); current = current.Add(f.cfg.IntervalDuration) {
		interval, err := f.NewBookingInterval(current, username)
		if err != nil {
			return nil, err
		}
		availableIntervals = append(availableIntervals, interval)
	}

	return availableIntervals, nil
}
