package sm

import (
	"fmt"

	"sm-instruction/internal/common/commonerrs"
)

type SkillType struct {
	s string
}

var (
	Engineering = SkillType{s: "Инженерные"}
	Researching = SkillType{s: "Исследовательские"}
	Social      = SkillType{s: "Социальные"}
	Creative    = SkillType{s: "Творческие"}
	Sportive    = SkillType{s: "Спортивные"}
)

var GeneralSkill = []SkillType{
	Engineering, Researching, Social,
}

var AdditionalSkill = []SkillType{
	Creative, Sportive,
}

var AllSkills = append(GeneralSkill, AdditionalSkill...)

func (s SkillType) String() string {
	return s.s
}

func (s SkillType) IsZero() bool {
	return s == SkillType{}
}

func NewSkillTypeFromString(s string) (SkillType, error) {
	switch s {
	case "Инженерные":
		return Engineering, nil
	case "Исследовательские":
		return Researching, nil
	case "Социальные":
		return Social, nil
	case "Творческие":
		return Creative, nil
	case "Спортивные":
		return Sportive, nil
	}
	return SkillType{}, commonerrs.NewInvalidInputError(
		fmt.Sprintf("invalid skill type %s, expected one of [%v]", s, AllSkills),
	)
}
