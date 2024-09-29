package query

import (
	"time"

	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

type User struct {
	Username string
	Role     string
}

type Slot struct {
	Start time.Time
	End   time.Time
	Whom  *string
}

type Grade struct {
	SkillType    string
	Points       int
	ActivityName string
	Time         time.Time
}

type Character struct {
	GroupName string
	Username  string
	Skills    map[string]int
	Rating    float64
	Slots     []Slot
	Grades    []Grade
	Start     *time.Time
	End       *time.Time
}

type Activity struct {
	Name        string
	FullName    string
	Description *string
	Location    *string
	Admins      []User
	Skills      []string
	MaxPoints   int
	Slots       []Slot
}

func convertUserToApp(u sm.User) User {
	return User{
		Username: u.Username,
		Role:     u.Role.String(),
	}
}

func convertUsersToApp(us []sm.User) []User {
	res := make([]User, len(us))
	for i, u := range us {
		res[i] = convertUserToApp(u)
	}
	return res
}

func convertSkillTypeToApp(s sm.SkillType) string {
	return s.String()
}

func convertSkillTypesToApp(ss []sm.SkillType) []string {
	res := make([]string, len(ss))
	for i, s := range ss {
		res[i] = convertSkillTypeToApp(s)
	}
	return res
}

func convertSkillsToApp(m map[sm.SkillType]int) map[string]int {
	res := make(map[string]int)
	for k, v := range m {
		res[convertSkillTypeToApp(k)] = v
	}
	return res
}

func convertSlotToApp(slot *sm.Slot) Slot {
	return Slot{
		Start: slot.Start,
		End:   slot.End,
		Whom:  slot.Whom,
	}
}

func convertSlotsToApp(slots []*sm.Slot) []Slot {
	res := make([]Slot, len(slots))
	for i, s := range slots {
		res[i] = convertSlotToApp(s)
	}
	return res
}

func convertGradeToApp(g sm.Grade) Grade {
	return Grade{
		SkillType:    g.SkillType.String(),
		Points:       g.Points,
		ActivityName: g.ActivityName,
		Time:         g.Time,
	}
}

func convertGradesToApp(gs []sm.Grade) []Grade {
	res := make([]Grade, len(gs))
	for i, g := range gs {
		res[i] = convertGradeToApp(g)
	}
	return res
}

func convertCharacterToApp(c *sm.Character) Character {
	return Character{
		Username:  c.Username,
		GroupName: c.GroupName,
		Skills:    convertSkillsToApp(c.Skills()),
		Rating:    c.Rating(),
		Slots:     convertSlotsToApp(c.Slots),
		Grades:    convertGradesToApp(c.Grades),
		Start:     c.StartedAt,
		End:       c.EndTime(),
	}
}

func convertCharactersToApp(cs []*sm.Character) []Character {
	res := make([]Character, len(cs))
	for i, c := range cs {
		res[i] = convertCharacterToApp(c)
	}
	return res
}

func convertActivityToApp(a *sm.Activity) Activity {
	return Activity{
		Name:        a.Name,
		FullName:    a.FullName,
		Description: a.Description,
		Location:    a.Location,
		Admins:      convertUsersToApp(a.Admins),
		Skills:      convertSkillTypesToApp(a.Skills),
		MaxPoints:   a.MaxPoints,
		Slots:       convertSlotsToApp(a.Slots),
	}
}

func convertActivitiesToApp(as []*sm.Activity) []Activity {
	res := make([]Activity, len(as))
	for i, a := range as {
		res[i] = convertActivityToApp(a)
	}
	return res
}
