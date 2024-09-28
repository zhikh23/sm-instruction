package sm

import (
	"github.com/zhikh23/sm-instruction/internal/common/commonerrs"
	"time"
)

type Grade struct {
	SkillType    SkillType
	Points       int
	ActivityName string
	Time         time.Time
}

func NewGrade(
	skillType SkillType,
	points int,
	activityName string,
	time time.Time,
) (Grade, error) {
	if skillType.IsZero() {
		return Grade{}, commonerrs.NewInvalidInputError("expected not empty skill type")
	}

	if points <= 0 {
		return Grade{}, commonerrs.NewInvalidInputError("expected positive number of points")
	}

	if activityName == "" {
		return Grade{}, commonerrs.NewInvalidInputError("expected non empty activity name")
	}

	if time.IsZero() {
		return Grade{}, commonerrs.NewInvalidInputError("expected non empty time")
	}

	return Grade{
		SkillType:    skillType,
		Points:       points,
		ActivityName: activityName,
		Time:         time,
	}, nil
}

func UnmarshallGradeFromDB(
	skillTypeStr string,
	points int,
	activityName string,
	time time.Time,
) (Grade, error) {
	if skillTypeStr == "" {
		return Grade{}, commonerrs.NewInvalidInputError("expected not empty skill type")
	}

	skillType, err := NewSkillTypeFromString(skillTypeStr)
	if err != nil {
		return Grade{}, err
	}

	if points <= 0 {
		return Grade{}, commonerrs.NewInvalidInputError("expected positive number of points")
	}

	if activityName == "" {
		return Grade{}, commonerrs.NewInvalidInputError("expected non empty activity name")
	}

	if time.IsZero() {
		return Grade{}, commonerrs.NewInvalidInputError("expected non empty time")
	}

	return Grade{
		SkillType:    skillType,
		Points:       points,
		ActivityName: activityName,
		Time:         time,
	}, nil
}

func (g Grade) IsZero() bool {
	return g == Grade{}
}
