package sm

import (
	"time"

	"sm-instruction/internal/common/commonerrs"
)

const BookInterval = 30 * time.Minute

type BookedTime struct {
	Username     string
	ActivityUUID string
	Start        time.Time
	Finish       time.Time
	CanBeRemoved bool
}

func NewBookedTime(
	username string,
	activityUUID string,
	start time.Time,
	canBeRemoved bool,
) (BookedTime, error) {
	if username == "" {
		return BookedTime{}, commonerrs.NewInvalidInputError("expected not empty username")
	}

	if activityUUID == "" {
		return BookedTime{}, commonerrs.NewInvalidInputError("expected not empty activity UUID")
	}

	if start.IsZero() {
		return BookedTime{}, commonerrs.NewInvalidInputError("expected non-zero time")
	}

	if !start.Round(BookInterval).Equal(start) {
		return BookedTime{}, commonerrs.NewInvalidInputErrorf(
			"invalid time %s; expected multiply of %s",
			start.String(), BookInterval.String(),
		)
	}

	finish := start.Add(BookInterval)

	return BookedTime{
		Username:     username,
		ActivityUUID: activityUUID,
		Start:        start,
		Finish:       finish,
		CanBeRemoved: canBeRemoved,
	}, nil
}

func MustNewBookedTime(
	username string,
	activityUUID string,
	start time.Time,
	canBeRemoved bool,
) BookedTime {
	b, err := NewBookedTime(username, activityUUID, start, canBeRemoved)
	if err != nil {
		panic(err)
	}
	return b
}

func (b BookedTime) IsZero() bool {
	return b == BookedTime{}
}

func (b BookedTime) TimeString() string {
	return b.Start.Format("15:04")
}
