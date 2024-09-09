package sm

import (
	"errors"
	"slices"
	"time"

	"sm-instruction/internal/common/commonerrs"
)

const AvailableSkillsNumber = 2

type Location struct {
	UUID            string
	Name            string
	Booked          []BookedTime
	Administrators  []User
	AvailableSkills []SkillType
}

func NewLocation(
	uuid string,
	name string,
	availableSkills []SkillType,
) (*Location, error) {
	if uuid == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty location uuid")
	}

	if name == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty location name")
	}

	if len(availableSkills) != AvailableSkillsNumber {
		return nil, commonerrs.NewInvalidInputErrorf(
			"invalid number of available skills: %d; expected %d",
			len(availableSkills), AvailableSkillsNumber,
		)
	}

	return &Location{
		UUID:            uuid,
		Name:            name,
		Booked:          make([]BookedTime, 0),
		Administrators:  make([]User, 0),
		AvailableSkills: availableSkills,
	}, nil
}

func MustNewLocation(
	uuid string,
	name string,
	availableSkills []SkillType,
) *Location {
	l, err := NewLocation(uuid, name, availableSkills)
	if err != nil {
		panic(err)
	}
	return l
}

func UnmarshallLocationFromDB(
	uuid string,
	name string,
	availableSkills []string,
	booked []BookedTime,
	administrators []User,
) (*Location, error) {
	if uuid == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty location uuid")
	}

	if name == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty location name")
	}

	if len(availableSkills) == 0 {
		return nil, commonerrs.NewInvalidInputError("expected not empty location available skills")
	}

	skills := make([]SkillType, len(availableSkills))
	for i, str := range availableSkills {
		skill, err := NewSkillTypeFromString(str)
		if err != nil {
			return nil, err
		}
		skills[i] = skill
	}

	if booked == nil {
		booked = make([]BookedTime, 0)
	}

	for _, interval := range booked {
		if interval.IsZero() {
			return nil, commonerrs.NewInvalidInputError("expected not empty booking interval")
		}
	}

	if administrators == nil {
		administrators = make([]User, 0)
	}

	for _, user := range administrators {
		if user.IsZero() {
			return nil, commonerrs.NewInvalidInputError("expected not empty administrators")
		}
	}

	return &Location{
		UUID:            uuid,
		Name:            name,
		Booked:          booked,
		Administrators:  administrators,
		AvailableSkills: skills,
	}, nil
}

func (l *Location) IsBooked(t time.Time) bool {
	for _, bt := range l.Booked {
		if bt.Time.Equal(t) {
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

func (l *Location) AddBooking(t time.Time, username string) error {
	if t.IsZero() {
		return commonerrs.NewInvalidInputError("expected not empty booking interval")
	}

	if err := l.CanBook(t); err != nil {
		return err
	}

	bookTime, err := NewBookedTime(t, username)
	if err != nil {
		return err
	}

	l.Booked = append(l.Booked, bookTime)

	return nil
}

var ErrLocationCannotIncSkill = errors.New("location cannot increment skill")

func (l *Location) Complete(char *Character, incSkill SkillType, score int) error {
	if err := char.RemoveBooking(l); err != nil {
		return err
	}

	if !slices.Contains(l.AvailableSkills, incSkill) {
		return ErrLocationCannotIncSkill
	}

	if err := char.IncSkill(incSkill, score); err != nil {
		return err
	}

	return nil
}

func (l *Location) AvailableTimes(to time.Time) []time.Time {
	times := make([]time.Time, 0)

	it := time.Now().Round(BookInterval)
	if it.Before(time.Now()) {
		it = it.Add(BookInterval)
	}
	times = append(times, it)

	for ; it.Before(to); it = it.Add(BookInterval) {
		if !l.IsBooked(it) {
			times = append(times, it)
		}
	}

	return times
}
