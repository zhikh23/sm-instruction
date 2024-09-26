package sm

import (
	"errors"
	"time"

	"sm-instruction/internal/common/commonerrs"
)

const TimeFormat = "15:04"

type Slot struct {
	Start time.Time
	End   time.Time
	Whom  *string
}

func NewSlot(
	start time.Time,
	end time.Time,
) (*Slot, error) {
	if start.IsZero() {
		return nil, commonerrs.NewInvalidInputError("expected not zero start time")
	}

	if !start.Round(time.Minute).Equal(start) {
		return nil, commonerrs.NewInvalidInputError("start time must be multiply of minute")
	}

	if end.IsZero() {
		return nil, commonerrs.NewInvalidInputError("expected not zero end time")
	}

	if !end.Round(time.Minute).Equal(end) {
		return nil, commonerrs.NewInvalidInputError("start time must be multiply of minute")
	}

	if !start.Before(end) {
		return nil, commonerrs.NewInvalidInputError("start time must be before end time")
	}

	return &Slot{
		Start: start,
		End:   end,
		Whom:  nil,
	}, nil
}

func MustNewSlot(
	start time.Time,
	end time.Time,
) *Slot {
	s, err := NewSlot(start, end)
	if err != nil {
		panic(err)
	}
	return s
}

func UnmarshallFromDB(
	start time.Time,
	end time.Time,
	whom *string,
) (*Slot, error) {
	if start.IsZero() {
		return nil, commonerrs.NewInvalidInputError("expected not zero start time")
	}

	if end.IsZero() {
		return nil, commonerrs.NewInvalidInputError("expected not zero end time")
	}

	if whom != nil && *whom == "" {
		return nil, commonerrs.NewInvalidInputError("expected not zero whom or nil")
	}

	return &Slot{
		Start: start,
		End:   end,
		Whom:  whom,
	}, nil
}

func (s *Slot) IsAvailable() bool {
	return s.Whom == nil
}

var ErrSlotHasAlreadyTaken = errors.New("slot has already taken")

func (s *Slot) Take(whom string) error {
	if !s.IsAvailable() {
		return ErrSlotHasAlreadyTaken
	}

	s.Whom = &whom

	return nil
}

var ErrSlotHasNotTaken = errors.New("slot has not taken")

func (s *Slot) Free() error {
	if s.IsAvailable() {
		return ErrSlotHasNotTaken
	}

	s.Whom = nil

	return nil
}

func SlotsIntersection(a, b []*Slot) []*Slot {
	res := make([]*Slot, 0, max(len(a), len(b)))
	for i := range a {
		for j := range b {
			if a[i].Start == b[j].Start {
				res = append(res, a[i])
			}
		}
	}
	return res
}
