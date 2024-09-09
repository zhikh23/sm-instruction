package query

import (
	"sm-instruction/internal/domain/sm"
	"time"
)

type User struct {
	ChatID   int64
	Username string
}

type Character struct {
	Head      User
	GroupName string
}

type BookedInterval struct {
	Time       time.Time
	ByUsername string
}

type Location struct {
	UUID            string
	Name            string
	Booked          []BookedInterval
	Administrators  []User
	AvailableSkills []string
}

func convertCharacterToApp(c *sm.Character) Character {
	return Character{
		Head:      convertUserToApp(c.Head),
		GroupName: c.GroupName,
	}
}

func convertUserToApp(u sm.User) User {
	return User{
		ChatID:   u.ChatID,
		Username: u.Username,
	}
}

func convertBookedIntervalToApp(i sm.BookedTime) BookedInterval {
	return BookedInterval{
		Time:       i.Time,
		ByUsername: i.ByUsername,
	}
}

func convertBookedIntervalsToApp(is []sm.BookedTime) []BookedInterval {
	res := make([]BookedInterval, len(is))
	for i, in := range is {
		res[i] = convertBookedIntervalToApp(in)
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

func convertLocationToApp(l *sm.Location) Location {
	return Location{
		UUID:            l.UUID,
		Name:            l.Name,
		Booked:          convertBookedIntervalsToApp(l.Booked),
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
