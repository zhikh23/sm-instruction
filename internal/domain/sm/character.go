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

type Character struct {
	Head      User
	GroupName string
	Skills    map[SkillType]int
	StartedAt *time.Time
	FinishAt  *time.Time

	BookedLocationTo   *time.Time
	BookedLocationUUID *string
}

func NewCharacter(
	head User,
	groupName string,
) (*Character, error) {
	if head.IsZero() {
		return nil, commonerrs.NewInvalidInputError("expected not empty head")
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
		Head:      head,
		GroupName: groupName,
		Skills:    skills,
		StartedAt: nil,
		FinishAt:  nil,

		BookedLocationTo:   nil,
		BookedLocationUUID: nil,
	}, nil
}

func MustNewCharacter(
	head User,
	groupName string,
) *Character {
	c, err := NewCharacter(head, groupName)
	if err != nil {
		panic(err)
	}
	return c
}

func NewCharacterFromDB(
	chatID int64,
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
	head, err := NewUser(chatID, username)
	if err != nil {
		return nil, err
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
		Head:      head,
		GroupName: groupName,
		Skills:    skills,
		StartedAt: startedAt,
		FinishAt:  finishAt,

		BookedLocationTo:   bookedLocationTo,
		BookedLocationUUID: bookedLocationUUID,
	}, nil
}

func ValidateGroupName(groupName string) error {
	pattern := `^[А-Я]{2}\d{1,2}\-\d{1,2}[А-Я]?$`
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

var ErrCharacterNotStarted = errors.New("character is not started")

func (c *Character) CheckStarted() error {
	if !c.IsStarted() {
		return ErrCharacterNotStarted
	}
	return nil
}

var ErrCharacterIsFinished = errors.New("character is finished")

func (c *Character) IsFinished() bool {
	if !c.IsStarted() {
		return false
	}
	return time.Now().After(*c.FinishAt)
}

func (c *Character) CheckFinished() error {
	if c.IsFinished() {
		return ErrCharacterIsFinished
	}
	return nil
}

func (c *Character) FinishTime() (time.Time, error) {
	if err := c.CheckFinished(); err != nil {
		return time.Time{}, err
	}

	return *c.FinishAt, nil
}

func (c *Character) Finish() error {
	if err := c.CheckFinished(); err != nil {
		return err
	}

	t := time.Now()
	c.FinishAt = &t

	return nil
}

func (c *Character) IsProcessing() bool {
	return c.IsStarted() && !c.IsFinished()
}

func (c *Character) Username() string {
	return c.Head.Username
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

func (c *Character) CanBook(loc *Location, interval BookingInterval) error {
	if !c.IsStarted() {
		return ErrCharacterNotStarted
	}

	if c.IsFinished() {
		return ErrCharacterIsFinished
	}

	if interval.From.After(*c.FinishAt) {
		return ErrCharacterBookingIsTooLate
	}

	if c.HasBooking() {
		return ErrCharacterAlreadyHasBooking
	}

	if err := loc.CheckBooked(interval); err != nil {
		return err
	}

	return nil
}

func (c *Character) Book(loc *Location, from time.Time, f BookingIntervalFactory) error {
	interval, err := f.NewBookingInterval(from, c.Head.Username)
	if err != nil {
		return err
	}

	if err := c.CanBook(loc, interval); err != nil {
		return err
	}

	if err := loc.AddBooking(interval); err != nil {
		return err
	}

	c.book(loc, interval)

	return nil
}

var ErrNotBooked = errors.New("not booked")
var ErrBookingLocationMismatch = errors.New("booking location mismatch")

func (c *Character) RemoveBooking(l *Location) error {
	if !c.HasBooking() {
		return ErrNotBooked
	}

	if *c.BookedLocationUUID != l.UUID {
		return ErrBookingLocationMismatch
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

func (c *Character) book(l *Location, i BookingInterval) {
	c.BookedLocationUUID = &l.UUID
	c.BookedLocationTo = &i.To
}
