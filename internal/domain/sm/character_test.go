package sm_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"sm-instruction/internal/common/commonerrs"
	"sm-instruction/internal/domain/sm"
)

const randomGroupName = "СМ1-11Б"
const randomUsername = "username"

func TestValidateGroupName(t *testing.T) {
	type testCase struct {
		GroupName string
		ExpectErr bool
	}

	cases := []testCase{
		{GroupName: "СМ1-11Б", ExpectErr: false},
		{GroupName: "СМ13-15Б", ExpectErr: false},
		{GroupName: "СМ5-15", ExpectErr: false},
		{GroupName: "СМ6-13", ExpectErr: false},
		{GroupName: "см1-11", ExpectErr: true},
		{GroupName: "ИУ7-34Б", ExpectErr: true},
		{GroupName: "СМ1-1", ExpectErr: true},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("group name: %s", tc.GroupName), func(t *testing.T) {
			err := sm.ValidateGroupName(tc.GroupName)
			if tc.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCharacter_Rating(t *testing.T) {
	t.Run("rating without factor", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, "СМ1-11Б")

		require.Zero(t, char.Rating())
		require.NoError(t, char.IncSkill(sm.Engineering, 3))
		require.NoError(t, char.IncSkill(sm.Researching, 4))
		require.Equal(t, 7.0, char.Rating())
	})

	t.Run("rating without general skills", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, "СМ1-11Б")

		require.Zero(t, char.Rating())
		require.NoError(t, char.IncSkill(sm.Creative, 3))
		require.NoError(t, char.IncSkill(sm.Sportive, 4))
		require.Equal(t, 0.0, char.Rating())
	})

	t.Run("rating with factor", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, "СМ1-11Б")

		require.Zero(t, char.Rating())
		require.NoError(t, char.IncSkill(sm.Engineering, 3))
		require.NoError(t, char.IncSkill(sm.Engineering, 4))
		require.NoError(t, char.IncSkill(sm.Creative, 3))
		require.NoError(t, char.IncSkill(sm.Creative, 3))
		require.NoError(t, char.IncSkill(sm.Sportive, 4))

		expected := (3 + 4) * (1 + 0.3 + 0.3 + 0.4)
		require.InDelta(t, expected, char.Rating(), 1e-3)
	})
}

func TestCharacter_IncSkill(t *testing.T) {
	t.Run("should increment character's skill", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, "СМ1-11Б")

		score := 5
		require.Zero(t, char.Skill(sm.Engineering))
		err := char.IncSkill(sm.Engineering, score)
		require.NoError(t, err)
		require.Equal(t, score, char.Skill(sm.Engineering))
	})

	t.Run("should return error if score is invalid", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, "СМ1-11Б")

		score := 10
		require.Zero(t, char.Skill(sm.Engineering))
		err := char.IncSkill(sm.Engineering, score)
		require.ErrorIs(t, err, sm.ErrInvalidScore)
		require.Zero(t, char.Skill(sm.Engineering))
	})
}

func TestCharacter_Booking(t *testing.T) {
	t.Run("should book location", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, randomGroupName)
		require.NoError(t, char.Start())

		loc := sm.MustNewLocation("1234", "345", "Description", "Test", []sm.SkillType{sm.Engineering, sm.Sportive})

		from := timeWithMinutes(30)
		err := char.Book(loc, from)
		require.NoError(t, err)
		require.True(t, char.HasBooking())
	})

	t.Run("should return error if character already has booking", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, randomGroupName)
		require.NoError(t, char.Start())

		loc := sm.MustNewLocation("1234", "345", "Description", "Test", []sm.SkillType{sm.Engineering, sm.Sportive})

		from := timeWithMinutes(30)
		err := char.Book(loc, from)
		require.NoError(t, err)

		from = timeWithMinutes(0)
		err = char.Book(loc, from)
		require.ErrorIs(t, err, sm.ErrCharacterAlreadyHasBooking)
	})

	t.Run("should return error if location is already booked", func(t *testing.T) {
		loc := createLocation()

		char1 := sm.MustNewCharacter("username2", randomGroupName)
		require.NoError(t, char1.Start())

		char2 := sm.MustNewCharacter("username1", randomGroupName)
		require.NoError(t, char2.Start())

		from := timeWithMinutes(30)
		err := char1.Book(loc, from)
		require.NoError(t, err)

		from = timeWithMinutes(30)
		err = char2.Book(loc, from)
		require.ErrorIs(t, err, sm.ErrLocationAlreadyBooked)
	})

	t.Run("should return error if book interval is invalid", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, randomGroupName)
		require.NoError(t, char.Start())

		loc := sm.MustNewLocation("1234", "345", "Desc", "Test", []sm.SkillType{sm.Engineering, sm.Sportive})

		from := timeWithMinutes(45)
		err := char.Book(loc, from)
		require.ErrorAs(t, err, &commonerrs.InvalidInputError{})
	})

	t.Run("should return error if booking is too late", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, randomGroupName)
		require.NoError(t, char.Start())

		loc := sm.MustNewLocation("1234", "345", "Desc", "Test", []sm.SkillType{sm.Engineering, sm.Sportive})

		from := time.Now().Add(-16 * time.Minute).Round(30 * time.Minute)
		err := char.Book(loc, from)
		require.ErrorIs(t, err, sm.ErrCharacterBookingIsTooClose)
	})

	t.Run("should return error if character has ended Instruction", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, randomGroupName)
		require.NoError(t, char.Start())

		loc := sm.MustNewLocation("1234", "345", "Desc", "Test", []sm.SkillType{sm.Engineering, sm.Sportive})

		from := time.Now().Add(sm.MaxDurationInstruction).Round(20 * time.Minute).Add(20 * time.Minute)
		err := char.Book(loc, from)
		require.ErrorIs(t, err, sm.ErrCharacterBookingIsTooLate)
	})
}

func TestCharacter_Start(t *testing.T) {
	char := sm.MustNewCharacter(randomUsername, randomGroupName)

	require.NoError(t, char.Start())
	require.True(t, char.IsStarted())
	require.True(t, char.IsProcessing())
}

func TestCharacter_Finish(t *testing.T) {
	t.Run("should finish started character", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, randomGroupName)

		require.NoError(t, char.Start())
		require.NoError(t, char.Finish())

		require.True(t, char.IsStarted())
		require.False(t, char.IsProcessing())
		require.True(t, char.IsFinished())
	})
}

func timeWithMinutes(minutes int) time.Time {
	now := time.Now()
	return time.Date(
		now.Year(), now.Month(), now.Day(),
		now.Hour(),
		minutes, 0, 0,
		time.Local,
	).Add(2 * time.Hour)
}

type createLocationParams struct {
	UUID        string
	Name        string
	Description string
	Where       string
	Skills      []sm.SkillType
}

type createLocationOption func(params *createLocationParams)

func createLocation(opts ...createLocationOption) *sm.Location {
	params := createLocationParams{
		UUID:        "1234",
		Name:        "Test",
		Description: "Test",
		Where:       "501м",
		Skills:      []sm.SkillType{sm.Engineering, sm.Sportive},
	}

	for _, opt := range opts {
		opt(&params)
	}

	return sm.MustNewLocation(params.UUID, params.Name, params.Description, params.Where, params.Skills)
}
