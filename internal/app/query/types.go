package query

import (
	"time"

	"sm-instruction/internal/domain/sm"
)

type User struct {
	Username string
	Role     string
}

type Character struct {
	Username  string
	GroupName string
	Skills    map[string]int
	StartedAt *time.Time
	FinishAt  *time.Time

	BookedLocationUUID *string
	BookedLocationTo   *time.Time
}

type BookedInterval struct {
	Time       time.Time
	ByUsername string
}

type Location struct {
	UUID            string
	Name            string
	Description     string
	Where           string
	Booked          []BookedInterval
	Administrators  []User
	AvailableSkills []string
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

func convertCharacterToApp(c *sm.Character) Character {
	return Character{
		Username:  c.Username,
		GroupName: c.GroupName,
		Skills:    convertSkillsToApp(c.Skills),
		StartedAt: c.StartedAt,
		FinishAt:  c.FinishAt,

		BookedLocationUUID: c.BookedLocationUUID,
		BookedLocationTo:   c.BookedLocationTo,
	}
}

func convertUserToApp(u sm.User) User {
	return User{
		Username: u.Username,
		Role:     u.Role.String(),
	}
}

func convertBookedTimeToApp(i sm.BookedTime) BookedInterval {
	return BookedInterval{
		Time:       i.Time,
		ByUsername: i.ByUsername,
	}
}

func convertBookedTimesToApp(is []sm.BookedTime) []BookedInterval {
	res := make([]BookedInterval, len(is))
	for i, in := range is {
		res[i] = convertBookedTimeToApp(in)
	}
	return res
}

func convertUsersToApp(us []sm.User) []User {
	res := make([]User, len(us))
	for i, u := range us {
		res[i] = convertUserToApp(u)
	}
	return res
}

func convertLocationToApp(l *sm.Location) Location {
	return Location{
		UUID:            l.UUID,
		Name:            l.Name,
		Description:     l.Description,
		Where:           l.Where,
		Booked:          convertBookedTimesToApp(l.Booked),
		Administrators:  convertUsersToApp(l.Administrators),
		AvailableSkills: convertSkillTypesToApp(l.AvailableSkills),
	}
}

func convertLocationsToApp(ls []*sm.Location) []Location {
	res := make([]Location, len(ls))
	for i, l := range ls {
		res[i] = convertLocationToApp(l)
	}
	return res
}
