package sm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"sm-instruction/internal/domain/sm"
)

func TestLocationFactory_NewLocation(t *testing.T) {
	t.Run("should create new location", func(t *testing.T) {
		l, err := sm.NewLocation("1234", "Test", []sm.SkillType{sm.Engineering, sm.Sportive})
		require.NoError(t, err)
		require.NotNil(t, l)
	})

	t.Run("should return error if skills number is invalid", func(t *testing.T) {
		l, err := sm.NewLocation("1234", "Test", []sm.SkillType{sm.Engineering, sm.Sportive, sm.Researching})
		require.Error(t, err)
		require.Nil(t, l)
	})
}

func TestLocation_AddBooking(t *testing.T) {
	t.Run("should add booking", func(t *testing.T) {
		loc := sm.MustNewLocation("1234", "Test", []sm.SkillType{sm.Researching, sm.Social})

		bt := timeWithMinutes(20)
		err := loc.AddBooking(bt, randomUsername)
		require.NoError(t, err)
		require.True(t, loc.IsBooked(bt))
		require.ErrorIs(t, loc.CanBook(bt), sm.ErrLocationAlreadyBooked)
	})

	t.Run("should return error if already booked", func(t *testing.T) {
		loc := sm.MustNewLocation("1234", "Test", []sm.SkillType{sm.Researching, sm.Social})

		bt := timeWithMinutes(20)
		require.NoError(t, loc.AddBooking(bt, randomUsername))
		err := loc.AddBooking(bt, randomUsername)
		require.ErrorIs(t, err, sm.ErrLocationAlreadyBooked)
	})
}

func TestLocation_Complete(t *testing.T) {
	t.Run("should complete character task", func(t *testing.T) {
		loc := sm.MustNewLocation("1234", "Test", []sm.SkillType{sm.Engineering, sm.Social})

		user := sm.MustNewUser(randomChatID, randomUsername)
		char := sm.MustNewCharacter(user, randomGroupName)
		require.NoError(t, char.Start())

		bt := timeWithMinutes(20)
		err := char.Book(loc, bt)
		require.NoError(t, err)

		score := 4
		require.Zero(t, char.Skill(sm.Engineering))
		err = loc.Complete(char, sm.Engineering, score)
		require.NoError(t, err)
		require.Equal(t, score, char.Skill(sm.Engineering))
	})

	t.Run("should return error if location cannot inc skill", func(t *testing.T) {
		loc := sm.MustNewLocation("1234", "Test", []sm.SkillType{sm.Researching, sm.Social})

		user := sm.MustNewUser(randomChatID, randomUsername)
		char := sm.MustNewCharacter(user, randomGroupName)
		require.NoError(t, char.Start())

		bt := timeWithMinutes(20)
		err := char.Book(loc, bt)
		require.NoError(t, err)

		score := 4
		require.Zero(t, char.Skill(sm.Engineering))
		err = loc.Complete(char, sm.Engineering, score)
		require.ErrorIs(t, err, sm.ErrLocationCannotIncSkill)
		require.Zero(t, char.Skill(sm.Engineering))
	})
}
