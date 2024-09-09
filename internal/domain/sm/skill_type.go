package sm

import (
	"fmt"

	"sm-instruction/internal/common/commonerrs"
)

type SkillType struct {
	s string
}

var (
	Engineering = SkillType{s: "engineering"}
	Researching = SkillType{s: "researching"}
	Social      = SkillType{s: "social"}
	Creative    = SkillType{s: "creative"}
	Sportive    = SkillType{s: "sportive"}
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
	case "engineering":
		return Engineering, nil
	case "researching":
		return Researching, nil
	case "social":
		return Social, nil
	case "creative":
		return Creative, nil
	case "sportive":
		return Sportive, nil
	}
	return SkillType{}, commonerrs.NewInvalidInputError(
		fmt.Sprintf("invalid skill type %s, expected one of [engineering, researching, creative, sportive]", s),
	)
}
