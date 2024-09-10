package sm

import (
	"time"

	"sm-instruction/internal/common/commonerrs"
)

const BookInterval = 30 * time.Minute

type BookedTime struct {
	Time       time.Time
	ByUsername string
}

func NewBookedTime(
	t time.Time,
	byUsername string,
) (BookedTime, error) {
	if t.IsZero() {
		return BookedTime{}, commonerrs.NewInvalidInputError("expected non-zero time")
	}

	if !t.Round(BookInterval).Equal(t) {
		return BookedTime{}, commonerrs.NewInvalidInputErrorf(
			"invalid time %s; expected multiply of %s",
			t.String(), BookInterval.String(),
		)
	}

	if byUsername == "" {
		return BookedTime{}, commonerrs.NewInvalidInputError("expected not empty username")
	}

	return BookedTime{
		Time:       t,
		ByUsername: byUsername,
	}, nil
}

func MustNewBookedTime(
	t time.Time,
	byUsername string,
) BookedTime {
	b, err := NewBookedTime(t, byUsername)
	if err != nil {
		panic(err)
	}
	return b
}

func (b BookedTime) IsZero() bool {
	return b == BookedTime{}
}

func (b BookedTime) TimeString() string {
	return b.Time.Format("15:04")
}
