package sm

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"sm-instruction/internal/common/commonerrs"
)

const MaxDurationInstruction = 4 * time.Hour
const MinimalDurationBeforeBooking = 5 * time.Minute

type Character struct {
	Username   string
	GroupName  string
	Skills     map[SkillType]int
	StartedAt  *time.Time
	FinishAt   *time.Time
	BookedTime *BookedTime
}

func NewCharacter(
	username string,
	groupName string,
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
	for _, s := range AllSkills {
		skills[s] = 0
	}

	return &Character{
		Username:   username,
		GroupName:  groupName,
		Skills:     skills,
		StartedAt:  nil,
		FinishAt:   nil,
		BookedTime: nil,
	}, nil
}

func MustNewCharacter(
	username string,
	groupName string,
) *Character {
	c, err := NewCharacter(username, groupName)
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
	startedAt *time.Time,
	finishAt *time.Time,
	bookedTime *BookedTime,
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

	return &Character{
		Username:   username,
		GroupName:  groupName,
		Skills:     skills,
		StartedAt:  startedAt,
		FinishAt:   finishAt,
		BookedTime: bookedTime,
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

var ErrCharacterAlreadyStarted = errors.New("character already started")

func (c *Character) Start() error {
	if c.IsStarted() {
		return ErrCharacterAlreadyStarted
	}

	start := time.Now()
	c.StartedAt = &start
	finish := start.Add(MaxDurationInstruction)
	c.FinishAt = &finish

	return nil
}

func (c *Character) IsStarted() bool {
	return c.StartedAt != nil
}

func (c *Character) IsFinished() bool {
	if !c.IsStarted() {
		return false
	}
	return time.Now().After(*c.FinishAt)
}

var ErrCharacterNotStarted = errors.New("character is not started")

func (c *Character) FinishTime() (time.Time, error) {
	if !c.IsStarted() {
		return time.Time{}, ErrCharacterNotStarted
	}

	return *c.FinishAt, nil
}

func (c *Character) IsProcessing() bool {
	return c.IsStarted() && !c.IsFinished()
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

func (c *Character) HasBooking() bool {
	return c.BookedTime != nil
}

var ErrCharacterHasNotBooking = errors.New("character has not booking")

func (c *Character) BookedTimeOrErr() (BookedTime, error) {
	if c.BookedTime == nil {
		return BookedTime{}, ErrCharacterHasNotBooking
	}
	return *c.BookedTime, nil
}

var ErrCharacterAlreadyHasBooking = errors.New("character already has booking")
var ErrCharacterBookingAfterFinish = errors.New("character booking after finish")
var ErrCharacterBookingIsTooClose = errors.New("character booking is too late")

func (c *Character) CanBook(t time.Time) error {
	if c.HasBooking() {
		return ErrCharacterAlreadyHasBooking
	}

	finish, err := c.FinishTime()
	if err != nil {
		return err
	}

	if t.After(finish) {
		return ErrCharacterBookingAfterFinish
	}

	if t.Add(-MinimalDurationBeforeBooking).Before(time.Now()) {
		return ErrCharacterBookingIsTooClose
	}

	return nil
}

func (c *Character) AddBooking(bookTime BookedTime) error {
	if err := c.CanBook(bookTime.Start); err != nil {
		return err
	}

	c.BookedTime = &bookTime

	return nil
}

func (c *Character) RemoveBooking() error {
	if !c.HasBooking() {
		return ErrCharacterHasNotBooking
	}

	c.BookedTime = nil

	return nil
}

func (c *Character) ratingBase() int {
	r := 0
	for _, s := range GeneralSkill {
		r += c.Skills[s]
	}
	return r
}

const lambda = 0.1

func (c *Character) ratingFactor() float64 {
	k := 1.0
	for _, s := range AdditionalSkill {
		k += lambda * float64(c.Skills[s])
	}
	return k
}
