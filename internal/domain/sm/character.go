package sm

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/zhikh23/sm-instruction/internal/common/commonerrs"
)

var ErrSlotAlreadyExists = errors.New("slot already exists")

type Character struct {
	GroupName string
	Username  string
	Skills    map[SkillType]int
	StartedAt *time.Time
	Slots     []*Slot
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

	times := make(map[time.Time]bool)
	for _, slot := range slots {
		if contains := times[slot.Start]; contains {
			return nil, ErrSlotAlreadyExists
		}
		times[slot.Start] = true
	}

	return &Character{
		GroupName: groupName,
		Username:  username,
		Skills:    skills,
		StartedAt: nil,
		Slots:     slots,
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
	groupName string,
	username string,
	skills map[SkillType]int,
	startedAt *time.Time,
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

	if slots == nil {
		slots = make([]*Slot, 0)
	}

	return &Character{
		Username:  username,
		GroupName: groupName,
		Skills:    skills,
		StartedAt: startedAt,
		Slots:     slots,
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

func (c *Character) IncSkill(t SkillType, score int) {
	c.Skills[t] += score
}

func (c *Character) AvailableSlots() []*Slot {
	return filter(c.Slots, slotIsAvailable)
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
	slot, ok := c.slotByTime(start)
	if !ok {
		return ErrSlotNotFound
	}

	return slot.Take(activityName)
}

func (c *Character) FreeSlot(start time.Time) error {
	slot, ok := c.slotByTime(start)
	if !ok {
		return ErrSlotNotFound
	}

	return slot.Free()
}

func (c *Character) slotByTime(start time.Time) (*Slot, bool) {
	for _, slot := range c.Slots {
		if slot.Start.Equal(start) {
			return slot, true
		}
	}
	return nil, false
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
