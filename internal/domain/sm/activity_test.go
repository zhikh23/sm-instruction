package sm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"sm-instruction/internal/domain/sm"
)

const randomActivityUUID = "1234"

func TestActivity_Award(t *testing.T) {
	t.Run("should award character", func(t *testing.T) {
		admin := sm.MustNewUser(randomUsername, sm.Administrator)
		act, err := sm.NewActivity(
			randomActivityUUID, "Викторина", []sm.User{admin}, []sm.SkillType{sm.Social}, 4, nil,
		)
		require.NoError(t, err)

		char := sm.MustNewCharacter(randomUsername, randomGroupName)

		err = act.Award(char, sm.Social, 4)
		require.NoError(t, err)
	})

	t.Run("should return error if activity can not award character", func(t *testing.T) {
		admin := sm.MustNewUser(randomUsername, sm.Administrator)
		act, err := sm.NewActivity(
			randomActivityUUID, "Викторина", []sm.User{admin}, []sm.SkillType{sm.Social}, 4, nil,
		)
		require.NoError(t, err)

		char := sm.MustNewCharacter(randomUsername, randomGroupName)

		err = act.Award(char, sm.Engineering, 4)
		require.ErrorIs(t, err, sm.ErrCannotIncSkill)
	})

	t.Run("should return error if activity awarding greater than max points", func(t *testing.T) {
		admin := sm.MustNewUser(randomUsername, sm.Administrator)
		act, err := sm.NewActivity(
			randomActivityUUID, "Викторина", []sm.User{admin}, []sm.SkillType{sm.Social}, 4, nil,
		)
		require.NoError(t, err)

		char := sm.MustNewCharacter(randomUsername, randomGroupName)

		err = act.Award(char, sm.Engineering, 5)
		require.ErrorIs(t, err, sm.ErrMaxPointsExceeded)
	})
}
