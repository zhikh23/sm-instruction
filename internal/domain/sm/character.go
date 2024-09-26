package sm

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"sm-instruction/internal/common/commonerrs"
)

type Character struct {
	GroupName string
	Username  string
	Skills    map[SkillType]int
	StartedAt *time.Time
	slots     map[time.Time]*Slot
}

func NewCharacter(
	groupName string,
	username string,
	slots []*Slot,
) (*Character, error) {
	if groupName == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty group")
	}

	if username == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty username")
	}

	if err := ValidateGroupName(groupName); err != nil {
		return nil, err
	}

	skills := make(map[SkillType]int)
	for _, s := range AllSkills {
		skills[s] = 0
	}

	mappedSlots, err := mapSlots(slots)
	if err != nil {
		return nil, err
	}

	return &Character{
		GroupName: groupName,
		Username:  username,
		Skills:    skills,
		StartedAt: nil,
		slots:     mappedSlots,
	}, nil
}

func MustNewCharacter(
	groupName string,
	username string,
	slots []*Slot,
) *Character {
	c, err := NewCharacter(groupName, username, slots)
	if err != nil {
		panic(err)
	}
	return c
}

func UnmarshallCharacterFromDB(
	username string,
	groupName string,
	engineeringSkill int,
	researchingSkill int,
	socialSkill int,
	creativeSkill int,
	sportiveSkill int,
	startedAt time.Time,
	slots []*Slot,
) (*Character, error) {
	if username == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty username")
	}

	if groupName == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty group")
	}

	if err := ValidateGroupName(groupName); err != nil {
		return nil, err
	}

	skills := make(map[SkillType]int)
	skills[Engineering] = engineeringSkill
	skills[Researching] = researchingSkill
	skills[Social] = socialSkill
	skills[Sportive] = sportiveSkill
	skills[Creative] = creativeSkill

	if slots == nil {
		slots = make([]*Slot, 0)
	}

	mappedSlots, err := mapSlots(slots)
	if err != nil {
		return nil, err
	}

	return &Character{
		Username:  username,
		GroupName: groupName,
		Skills:    skills,
		StartedAt: &startedAt,
		slots:     mappedSlots,
	}, nil
}

func ValidateGroupName(groupName string) error {
	pattern := `^СМ\d{1,2}\-\d{2,3}[Б]?$`
	if !regexp.MustCompile(pattern).MatchString(groupName) {
		return commonerrs.NewInvalidInputError(
			fmt.Sprintf(
				"invalid group name %s; expected match regular expression %s",
				groupName, pattern,
			))
	}
	return nil
}

func (c *Character) Skill(t SkillType) int {
	return c.Skills[t]
}

func (c *Character) Rating() float64 {
	return c.ratingFactor() * float64(c.ratingBase())
}

func (c *Character) Slots() []*Slot {
	slots := make([]*Slot, 0, len(c.slots))
	for _, slot := range c.slots {
		slots = append(slots, slot)
	}
	return slots
}

func (c *Character) IncSkill(t SkillType, score int) {
	c.Skills[t] += score
}

func (c *Character) AvailableSlots() []*Slot {
	return filter(c.Slots(), slotIsAvailable)
}

var ErrAlreadyStarted = errors.New("character already started")

func (c *Character) Start() error {
	if c.StartedAt != nil {
		return ErrAlreadyStarted
	}

	t := time.Now()
	c.StartedAt = &t

	return nil
}

var ErrSlotNotFound = errors.New("slot not found")

func (c *Character) TakeSlot(start time.Time, activityName string) error {
	fmt.Println(c.slots)
	slot, ok := c.slots[start]
	if !ok {
		return ErrSlotNotFound
	}

	return slot.Take(activityName)
}

func (c *Character) FreeSlot(start time.Time) error {
	slot, ok := c.slots[start]
	if !ok {
		return ErrSlotNotFound
	}

	return slot.Free()
}

func mapSlots(slots []*Slot) (map[time.Time]*Slot, error) {
	m := make(map[time.Time]*Slot)
	for _, slot := range slots {
		if _, ok := m[slot.Start]; ok {
			return nil, ErrSlotAlreadyExists
		}
		m[slot.Start] = slot
	}
	return m, nil
}

func (c *Character) ratingBase() int {
	r := 0
	for _, s := range GeneralSkill {
		r += c.Skills[s]
	}
	return r
}

const lambda = 0.25

func (c *Character) ratingFactor() float64 {
	k := 1.0
	for _, s := range AdditionalSkill {
		k += lambda * float64(c.Skills[s])
	}
	return k
}

func slotIsAvailable(slot *Slot) bool {
	return slot.IsAvailable()
}

func filter[T any](collection []T, predicate func(T) bool) []T {
	res := make([]T, 0, len(collection))
	for _, x := range collection {
		if predicate(x) {
			res = append(res, x)
		}
	}
	return res
}
