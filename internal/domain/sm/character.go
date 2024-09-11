package sm

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"sm-instruction/internal/common/commonerrs"
)

const MaxScore = 5
const MaxDurationInstruction = 4 * time.Hour
const MinimalDurationBeforeBooking = 5 * time.Minute

type Character struct {
	Username  string
	GroupName string
	Skills    map[SkillType]int
	StartedAt *time.Time
	FinishAt  *time.Time

	BookedLocationTo   *time.Time
	BookedLocationUUID *string
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
		Username:  username,
		GroupName: groupName,
		Skills:    skills,
		StartedAt: nil,
		FinishAt:  nil,

		BookedLocationTo:   nil,
		BookedLocationUUID: nil,
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
	bookedLocationTo *time.Time,
	bookedLocationUUID *string,
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
		Username:  username,
		GroupName: groupName,
		Skills:    skills,
		StartedAt: startedAt,
		FinishAt:  finishAt,

		BookedLocationTo:   bookedLocationTo,
		BookedLocationUUID: bookedLocationUUID,
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

var ErrCharacterAlreadyFinished = errors.New("character already finished")

func (c *Character) Finish() error {
	if !c.IsStarted() {
		return ErrCharacterNotStarted
	}

	if c.IsFinished() {
		return ErrCharacterAlreadyFinished
	}

	t := time.Now()
	c.FinishAt = &t

	return nil
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

var ErrInvalidScore = errors.New("invalid score")

func (c *Character) IncSkill(t SkillType, score int) error {
	if score < 0 || score > MaxScore {
		return ErrInvalidScore
	}

	c.Skills[t] += score

	return nil
}

func (c *Character) HasBooking() bool {
	if c.BookedLocationUUID == nil || c.BookedLocationTo == nil {
		return false
	}
	return time.Now().Before(*c.BookedLocationTo)
}

var ErrCharacterAlreadyHasBooking = errors.New("character already has booking")
var ErrCharacterBookingIsTooLate = errors.New("booking is too late")
var ErrCharacterBookingIsTooClose = errors.New("booking is too close")

func (c *Character) CanBook(loc *Location, t time.Time) error {
	if !c.IsStarted() {
		return ErrCharacterNotStarted
	}

	if c.IsFinished() {
		return ErrCharacterAlreadyFinished
	}

	if t.After(*c.FinishAt) {
		return ErrCharacterBookingIsTooLate
	}

	if t.Sub(time.Now()) < MinimalDurationBeforeBooking {
		return ErrCharacterBookingIsTooClose
	}

	if c.HasBooking() {
		return ErrCharacterAlreadyHasBooking
	}

	if err := loc.CanBook(t); err != nil {
		return err
	}

	return nil
}

func (c *Character) Book(loc *Location, t time.Time) error {
	if err := c.CanBook(loc, t); err != nil {
		return err
	}

	if err := loc.AddBooking(t, c.Username); err != nil {
		return err
	}

	c.book(loc, t)

	return nil
}

var ErrCharacterHasNotBooking = errors.New("character has not booked location")
var ErrCharacterBookingMismatch = errors.New("character booking location mismatch")

func (c *Character) BookedLocation() (string, error) {
	if !c.HasBooking() {
		return "", ErrCharacterHasNotBooking
	}

	return *c.BookedLocationUUID, nil
}

func (c *Character) RemoveBooking(l *Location) error {
	if !c.HasBooking() {
		return ErrCharacterHasNotBooking
	}

	if *c.BookedLocationUUID != l.UUID {
		return ErrCharacterBookingMismatch
	}

	c.BookedLocationUUID = nil
	c.BookedLocationTo = nil

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

func (c *Character) book(l *Location, t time.Time) {
	c.BookedLocationUUID = &l.UUID
	t = t.Add(BookInterval)
	c.BookedLocationTo = &t
}
