package sm

import (
	"errors"
	"slices"

	"sm-instruction/internal/common/commonerrs"
)

const AvailableSkillsNumber = 2

type Location struct {
	UUID            string
	Name            string
	Booked          []BookingInterval
	Administrators  []User
	AvailableSkills []SkillType
}

func UnmarshallLocationFromDB(
	uuid string,
	name string,
	availableSkills []string,
	booked []BookingInterval,
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
		booked = make([]BookingInterval, 0)
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

func (l *Location) IsBooked(i BookingInterval) bool {
	for _, b := range l.Booked {
		if b.IsIntersects(i) {
			return true
		}
	}
	return false
}

var ErrLocationIntervalHasAlreadyBooked = errors.New("interval already booked")

func (l *Location) CheckBooked(i BookingInterval) error {
	if l.IsBooked(i) {
		return ErrLocationIntervalHasAlreadyBooked
	}
	return nil
}

func (l *Location) AddBooking(i BookingInterval) error {
	if i.IsZero() {
		return commonerrs.NewInvalidInputError("expected not empty booking interval")
	}

	if err := l.CheckBooked(i); err != nil {
		return err
	}

	l.Booked = append(l.Booked, i)

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
		Name:            name,
		Booked:          make([]BookingInterval, 0),
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
