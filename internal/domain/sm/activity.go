package sm

import (
	"errors"
	"slices"
	"sm-instruction/pkg/funcs"
	"time"

	"sm-instruction/internal/common/commonerrs"
)

type Activity struct {
	Name        string
	Description *string
	Location    *string
	Admins      []User
	Skills      []SkillType
	MaxPoints   int
	slots       map[time.Time]*Slot
}

func NewActivity(
	name string,
	description *string,
	location *string,
	admins []User,
	skills []SkillType,
	maxPoints int,
	slots []*Slot,
) (*Activity, error) {
	if name == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty name")
	}

	if description != nil && *description == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty description or nil")
	}

	if location != nil && *location == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty name or nil")
	}

	if len(admins) == 0 {
		return nil, commonerrs.NewInvalidInputError("expected not empty admins")
	}

	for _, admin := range admins {
		if admin.Role != Administrator {
			return nil, commonerrs.NewInvalidInputErrorf(
				"expected user has role %q, got %q", Administrator.String(), admin.Role.String(),
			)
		}
	}

	if len(skills) == 0 {
		return nil, commonerrs.NewInvalidInputError("expected not empty skills")
	}

	for _, skill := range skills {
		if skill.IsZero() {
			return nil, commonerrs.NewInvalidInputError("expected not empty skill")
		}
	}

	if maxPoints <= 0 {
		return nil, commonerrs.NewInvalidInputError("expected positive max points")
	}

	if slots == nil {
		slots = make([]*Slot, 0)
	}

	mappedSlots, err := mapSlots(slots)
	if err != nil {
		return nil, err
	}

	return &Activity{
		Name:        name,
		Description: description,
		Location:    location,
		Admins:      admins,
		Skills:      skills,
		MaxPoints:   maxPoints,
		slots:       mappedSlots,
	}, nil
}

func MustNewActivity(
	name string,
	description *string,
	location *string,
	admins []User,
	skills []SkillType,
	maxPoints int,
	slots []*Slot,
) *Activity {
	a, err := NewActivity(name, description, location, admins, skills, maxPoints, slots)
	if err != nil {
		panic(err)
	}
	return a
}

func UnmarshallActivityFromDB(
	uuid string,
	name string,
	description *string,
	location *string,
	admins []User,
	skillsStr []string,
	maxPoints int,
	slots []*Slot,
) (*Activity, error) {
	if uuid == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty uuid")
	}

	if name == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty name")
	}

	if description != nil && *description == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty description or nil")
	}

	if location != nil && *location == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty name or nil")
	}

	if len(admins) == 0 {
		return nil, commonerrs.NewInvalidInputError("expected not empty admins")
	}

	if len(skillsStr) == 0 {
		return nil, commonerrs.NewInvalidInputError("expected not empty skills")
	}

	if maxPoints == 0 {
		return nil, commonerrs.NewInvalidInputError("expected not empty max points")
	}

	skills := make([]SkillType, len(skillsStr))
	for i, str := range skillsStr {
		skill, err := NewSkillTypeFromString(str)
		if err != nil {
			return nil, err
		}
		skills[i] = skill
	}

	mappedSlots, err := mapSlots(slots)
	if err != nil {
		return nil, err
	}

	return &Activity{
		Name:      name,
		Admins:    admins,
		Skills:    skills,
		MaxPoints: maxPoints,
		slots:     mappedSlots,
	}, nil
}

func (a *Activity) Slots() []*Slot {
	slots := make([]*Slot, 0, len(a.slots))
	for _, s := range a.slots {
		slots = append(slots, s)
	}
	return slots
}

var ErrCannotIncSkill = errors.New("cannot increment skill")
var ErrMaxPointsExceeded = errors.New("max points exceeded")

func (a *Activity) Award(char *Character, skill SkillType, points int) error {
	if points > a.MaxPoints {
		return ErrMaxPointsExceeded
	}

	if !slices.Contains(a.Skills, skill) {
		return ErrCannotIncSkill
	}

	char.IncSkill(skill, points)

	return nil
}

func (a *Activity) TakeSlot(start time.Time, groupName string) error {
	slot, ok := a.slots[start]
	if !ok {
		return ErrSlotNotFound
	}

	return slot.Take(groupName)
}

func (a *Activity) FreeSlot(start time.Time) error {
	slot, ok := a.slots[start]
	if !ok {
		return ErrSlotNotFound
	}

	return slot.Free()
}

func (a *Activity) AvailableSlots() []*Slot {
	return funcs.Filter(a.Slots(), func(slot *Slot) bool {
		return slot.IsAvailable()
	})
}
