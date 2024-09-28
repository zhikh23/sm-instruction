package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/zhikh23/sm-instruction/internal/adapters"

	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

const slotDuration = 20 * time.Minute

func main() {
	ctx := context.Background()

	activitiesProvider := adapters.NewDefaultGSActivitiesProvider()
	activities, err := activitiesProvider.Activities(ctx)
	if err != nil {
		log.Fatal(err)
	}

	charactersProvider := adapters.NewDefaultGSCharactersProvider()
	templateCharacters, err := charactersProvider.Characters(ctx)
	if err != nil {
		log.Fatal(err)
	}

	mapGroupToUsername := make(map[string]string)
	for _, char := range templateCharacters {
		mapGroupToUsername[char.GroupName] = char.Username
	}

	usersRepos, closeUsers := adapters.NewPGUsersRepository()
	defer func() {
		_ = closeUsers()
	}()

	charsRepos, closeChars := adapters.NewPGCharactersRepository()
	defer func() {
		_ = closeChars()
	}()

	activitiesRepos, closeActivities := adapters.NewPGActivitiesRepository()
	defer func() {
		_ = closeActivities()
	}()

	groups := make(map[string]bool)
	for _, act := range activities {
		for _, slot := range act.Slots() {
			if !slot.IsAvailable() {
				groups[*slot.Whom] = true
			}
		}
	}

	times := slotTimes()

	users := make(map[string]sm.User)
	for _, act := range activities {
		for _, admin := range act.Admins {
			users[admin.Username] = admin
		}
	}

	chars := make(map[string]*sm.Character, len(groups))
	for group := range groups {
		username, ok := mapGroupToUsername[group]
		if !ok {
			log.Fatalf("GroupName for group %s not found", group)
		}
		user := sm.MustNewUser(username, sm.Participant)
		users[user.Username] = user
		char := sm.MustNewCharacter(group, username, emptySlots(times))
		chars[char.GroupName] = char
	}

	for _, user := range users {
		err = usersRepos.Save(ctx, user)
		if err != nil && !errors.Is(err, sm.ErrUserAlreadyExists) {
			log.Fatalf("Failed to save user %s: %s", user.Username, err.Error())
		}
	}

	for _, char := range chars {
		err = charsRepos.Save(ctx, char)
		if err != nil && !errors.Is(err, sm.ErrCharacterAlreadyExists) {
			log.Fatalf("Failed to save character %s: %s", char.GroupName, err.Error())
		}
	}

	for _, act := range activities {
		err = activitiesRepos.Save(ctx, act)
		if err != nil && !errors.Is(err, sm.ErrActivityAlreadyExists) {
			log.Fatalf("Failed to save actitvity %q: %s\n", act.Name, err.Error())
		}
	}

	for _, act := range activities {
		for _, slot := range act.Slots() {
			if slot.IsAvailable() {
				continue
			}
			group := *slot.Whom
			if err = charsRepos.Update(ctx, group, func(innerCtx context.Context, char *sm.Character) error {
				return char.TakeSlot(slot.Start, act.Name)
			}); err != nil {
				log.Fatalf("Failed to update char %s: %s", group, err.Error())
			}
		}
	}
}

func slotTimes() []time.Time {
	times := make([]time.Time, 0)
	first := todayTime(11, 20)
	last := todayTime(17, 20)
	for it := first; !it.After(last); it = it.Add(slotDuration) {
		times = append(times, it)
	}
	return times
}

func todayTime(hours int, minutes int) time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hours, minutes, 0, 0, time.Local)
}

func emptySlots(times []time.Time) []*sm.Slot {
	slots := make([]*sm.Slot, len(times))
	for i, t := range times {
		slots[i] = sm.MustNewSlot(t, t.Add(slotDuration))
	}
	return slots
}
