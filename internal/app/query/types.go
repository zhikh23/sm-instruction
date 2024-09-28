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

type Character struct {
	GroupName string
	Username  string
	Skills    map[string]int
	Rating    float64
	Slots     []Slot
}

type Activity struct {
	Name      string
	Admins    []User
	Skills    []string
	MaxPoints int
	Slots     []Slot
}

type ActivityWithLocation struct {
	Activity
	Description string
	Location    string
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

func convertCharacterToApp(c *sm.Character) Character {
	return Character{
		Username:  c.Username,
		GroupName: c.GroupName,
		Skills:    convertSkillsToApp(c.Skills),
		Rating:    c.Rating(),
		Slots:     convertSlotsToApp(c.Slots()),
	}
}

func convertActivityToApp(a *sm.Activity) Activity {
	return Activity{
		Name:      a.Name,
		Admins:    convertUsersToApp(a.Admins),
		Skills:    convertSkillTypesToApp(a.Skills),
		MaxPoints: a.MaxPoints,
		Slots:     convertSlotsToApp(a.Slots()),
	}
}

func convertActivityWithLocation(a *sm.Activity) ActivityWithLocation {
	return ActivityWithLocation{
		Activity:    convertActivityToApp(a),
		Description: *a.Description,
		Location:    *a.Location,
	}
}

func convertActivitiesWithLocations(a []*sm.Activity) []ActivityWithLocation {
	res := make([]ActivityWithLocation, len(a))
	for i, act := range a {
		res[i] = convertActivityWithLocation(act)
	}
	return res
}
