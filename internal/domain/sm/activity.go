package sm

import (
	"errors"
	"slices"

	"sm-instruction/internal/common/commonerrs"
)

type Activity struct {
	UUID      string
	Name      string
	Admins    []User
	Skills    []SkillType
	MaxPoints int
	Location  *Location
}

func NewActivity(
	uuid string,
	name string,
	admins []User,
	skills []SkillType,
	maxPoints int,
	location *Location,
) (*Activity, error) {
	if uuid == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty uuid")
	}

	if name == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty name")
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

	return &Activity{
		UUID:      uuid,
		Name:      name,
		Admins:    admins,
		Skills:    skills,
		MaxPoints: maxPoints,
		Location:  location,
	}, nil
}

func MustNewActivity(
	uuid string,
	name string,
	admins []User,
	skills []SkillType,
	maxPoints int,
	location *Location,
) *Activity {
	a, err := NewActivity(uuid, name, admins, skills, maxPoints, location)
	if err != nil {
		panic(err)
	}
	return a
}

func UnmarshallActivityFromDB(
	uuid string,
	name string,
	admins []User,
	skillsStr []string,
	maxPoints int,
) (*Activity, error) {
	if uuid == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty uuid")
	}

	if name == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty name")
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

	return &Activity{
		UUID:      uuid,
		Name:      name,
		Admins:    admins,
		Skills:    skills,
		MaxPoints: maxPoints,
	}, nil
}

var ErrActivityHasNotLocation = errors.New("activity has not Location")

func (a *Activity) LocationOrErr() (*Location, error) {
	if a.Location == nil {
		return nil, ErrActivityHasNotLocation
	}

	return a.Location, nil
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
