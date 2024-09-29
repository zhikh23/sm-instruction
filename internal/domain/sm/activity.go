package sm

import (
	"errors"
	"slices"
	"time"

	"github.com/zhikh23/sm-instruction/internal/common/commonerrs"
	"github.com/zhikh23/sm-instruction/pkg/funcs"
)

type Activity struct {
	Name        string
	FullName    string
	Description *string
	Location    *string
	Admins      []User
	Skills      []SkillType
	MaxPoints   int
	Slots       []*Slot
}

func NewActivity(
	name string,
	fullName string,
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

	if fullName == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty full name")
	}

	if description != nil && *description == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty description or nil")
	}

	if location != nil && *location == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty name or nil")
	}

	if admins == nil {
		admins = make([]User, 0)
	}

	for _, admin := range admins {
		if admin.Role != Administrator {
			return nil, commonerrs.NewInvalidInputErrorf(
				"expected user has role %q, got %q", Administrator.String(), admin.Role.String(),
			)
		}
	}

	if skills == nil {
		skills = make([]SkillType, 0)
	}

	for _, skill := range skills {
		if skill.IsZero() {
			return nil, commonerrs.NewInvalidInputError("expected not empty skill")
		}
	}

	if maxPoints < 0 {
		return nil, commonerrs.NewInvalidInputError("expected non-negative max points")
	}

	if slots == nil {
		slots = make([]*Slot, 0)
	}

	times := make(map[time.Time]bool)
	for _, slot := range slots {
		if contains := times[slot.Start]; contains {
			return nil, ErrSlotAlreadyExists
		}
		times[slot.Start] = true
	}

	return &Activity{
		Name:        name,
		FullName:    fullName,
		Description: description,
		Location:    location,
		Admins:      admins,
		Skills:      skills,
		MaxPoints:   maxPoints,
		Slots:       slots,
	}, nil
}

func UnmarshallActivityFromDB(
	name string,
	fullName string,
	description *string,
	location *string,
	admins []User,
	skillsStr []string,
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

	if admins == nil {
		admins = make([]User, 0)
	}

	skills := make([]SkillType, len(skillsStr))
	for i, str := range skillsStr {
		skill, err := NewSkillTypeFromString(str)
		if err != nil {
			return nil, err
		}
		skills[i] = skill
	}

	return &Activity{
		Name:        name,
		FullName:    fullName,
		Description: description,
		Location:    location,
		Admins:      admins,
		Skills:      skills,
		MaxPoints:   maxPoints,
		Slots:       slots,
	}, nil
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

	return char.GiveGrade(skill, points, a.Name)
}

func (a *Activity) TakeSlot(start time.Time, groupName string) error {
	slot, ok := a.slotByTime(start)
	if !ok {
		return ErrSlotNotFound
	}

	return slot.Take(groupName)
}

func (a *Activity) FreeSlot(start time.Time) error {
	slot, ok := a.slotByTime(start)
	if !ok {
		return ErrSlotNotFound
	}

	return slot.Free()
}

func (a *Activity) AvailableSlots() []*Slot {
	return funcs.Filter(a.Slots, func(slot *Slot) bool {
		return slot.IsAvailable()
	})
}

func (a *Activity) HasTaken(groupName string) bool {
	for _, slot := range a.Slots {
		if !slot.IsAvailable() && *slot.Whom == groupName {
			return true
		}
	}
	return false
}

func (a *Activity) slotByTime(start time.Time) (*Slot, bool) {
	for _, slot := range a.Slots {
		if slot.Start.Equal(start) {
			return slot, true
		}
	}
	return nil, false
}
