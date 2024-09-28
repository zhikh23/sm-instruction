package sm

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"time"

	"github.com/zhikh23/sm-instruction/internal/common/commonerrs"

	"github.com/zhikh23/sm-instruction/pkg/funcs"
)

const InstructionDuration = 4 * 4 * time.Hour // TODO!

var ErrSlotAlreadyExists = errors.New("slot already exists")

type Character struct {
	GroupName string
	Username  string
	StartedAt *time.Time
	Slots     []*Slot
	Grades    []Grade
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
		StartedAt: nil,
		Slots:     slots,
		Grades:    make([]Grade, 0),
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
	startedAt *time.Time,
	slots []*Slot,
	grades []Grade,
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

	if grades == nil {
		grades = make([]Grade, 0)
	}

	return &Character{
		Username:  username,
		GroupName: groupName,
		StartedAt: startedAt,
		Slots:     slots,
		Grades:    grades,
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

func (c *Character) Rating() float64 {
	return c.ratingFactor() * float64(c.ratingBase())
}

func (c *Character) GiveGrade(skillType SkillType, points int, activityName string) error {
	grade, err := NewGrade(skillType, points, activityName, time.Now())
	if err != nil {
		return err
	}

	c.Grades = append(c.Grades, grade)

	return nil
}

func (c *Character) AvailableSlots() []*Slot {
	return filter(c.Slots, slotIsAvailable)
}

func (c *Character) Start() error {
	t := time.Now()
	c.StartedAt = &t

	// Обрезаем слоты до нужного промежутка времени.
	c.Slots = funcs.Filter(c.Slots, func(slot *Slot) bool {
		return slot.Start.Before(c.StartedAt.Add(InstructionDuration))
	})

	return nil
}

func (c *Character) IsStarted() bool {
	return c.StartedAt != nil
}

func (c *Character) EndTime() *time.Time {
	if !c.IsStarted() {
		return nil
	}

	v := c.StartedAt.Add(InstructionDuration)

	return &v
}

var ErrSlotNotFound = errors.New("slot not found")
var ErrSlotIsTooLate = errors.New("slot is too late")

func (c *Character) TakeSlot(start time.Time, activityName string) error {
	if c.IsStarted() && start.After(*c.EndTime()) {
		return ErrSlotIsTooLate
	}

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

func (c *Character) Skills() map[SkillType]int {
	skills := make(map[SkillType]int)
	for _, skill := range AllSkills {
		skills[skill] = c.sumPoints(func(st SkillType) bool {
			return st == skill
		})
	}
	return skills
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
	isGeneral := func(st SkillType) bool {
		return slices.Contains(GeneralSkill, st)
	}
	return c.sumPoints(isGeneral)
}

const lambda = 1.0 / 72 // Не спрашивайте, почему.

func (c *Character) ratingFactor() float64 {
	isAdditional := func(st SkillType) bool {
		return slices.Contains(AdditionalSkill, st)
	}
	return 1.0 + lambda*float64(c.sumPoints(isAdditional))
}

func (c *Character) sumPoints(predicate func(st SkillType) bool) int {
	r := 0
	for _, g := range c.Grades {
		if predicate(g.SkillType) {
			r += g.Points
		}
	}
	return r
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
