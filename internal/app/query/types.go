package query

import (
	"time"

	"sm-instruction/internal/domain/sm"
)

type User struct {
	Username string
	Role     string
}

type BookedTime struct {
	Username     string
	ActivityUUID string
	Start        time.Time
	Finish       time.Time
	CanBeRemoved bool
}

type Character struct {
	Username   string
	GroupName  string
	Skills     map[string]int
	StartedAt  *time.Time
	FinishAt   *time.Time
	BookedTime *BookedTime
}

type Activity struct {
	UUID      string
	Name      string
	Admins    []User
	Skills    []string
	MaxPoints int
}

type Location struct {
	UUID        string
	Name        string
	Admins      []User
	Skills      []string
	MaxPoints   int
	Description string
	Where       string
	BookedTimes []BookedTime
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

func convertBookedTimeToApp(bt sm.BookedTime) BookedTime {
	return BookedTime{
		Username:     bt.Username,
		ActivityUUID: bt.ActivityUUID,
		Start:        bt.Start,
		Finish:       bt.Finish,
		CanBeRemoved: bt.CanBeRemoved,
	}
}

func convertBookedTimeToAppOrNil(bt *sm.BookedTime) *BookedTime {
	if bt == nil {
		return nil
	}
	res := convertBookedTimeToApp(*bt)
	return &res
}

func convertBookedTimesToApp(bs []sm.BookedTime) []BookedTime {
	res := make([]BookedTime, len(bs))
	for i, b := range bs {
		res[i] = convertBookedTimeToApp(b)
	}
	return res
}

func convertCharacterToApp(c *sm.Character) Character {
	return Character{
		Username:   c.Username,
		GroupName:  c.GroupName,
		Skills:     convertSkillsToApp(c.Skills),
		StartedAt:  c.StartedAt,
		FinishAt:   c.FinishAt,
		BookedTime: convertBookedTimeToAppOrNil(c.BookedTime),
	}
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

func convertActivityToApp(a *sm.Activity) Activity {
	return Activity{
		UUID:      a.UUID,
		Name:      a.Name,
		Admins:    convertUsersToApp(a.Admins),
		Skills:    convertSkillTypesToApp(a.Skills),
		MaxPoints: a.MaxPoints,
	}
}

func convertActivitiesToApp(as []*sm.Activity) []Activity {
	res := make([]Activity, len(as))
	for i, a := range as {
		res[i] = convertActivityToApp(a)
	}
	return res
}

func convertLocationToApp(a *sm.Activity) Location {
	if a.Location == nil {
		panic("activity has not location")
	}
	return Location{
		UUID:        a.UUID,
		Name:        a.Name,
		Admins:      convertUsersToApp(a.Admins),
		Skills:      convertSkillTypesToApp(a.Skills),
		MaxPoints:   a.MaxPoints,
		Description: a.Location.Description,
		Where:       a.Location.Where,
		BookedTimes: convertBookedTimesToApp(a.Location.BookedTimes),
	}
}

func convertLocationsToApp(as []*sm.Activity) []Location {
	res := make([]Location, len(as))
	for i, a := range as {
		res[i] = convertLocationToApp(a)
	}
	return res
}
