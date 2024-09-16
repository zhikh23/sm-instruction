package sm_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

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
		char.IncSkill(sm.Engineering, 3)
		char.IncSkill(sm.Researching, 4)
		require.Equal(t, 7.0, char.Rating())
	})

	t.Run("rating without general skills", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, "СМ1-11Б")

		require.Zero(t, char.Rating())
		char.IncSkill(sm.Creative, 3)
		char.IncSkill(sm.Sportive, 4)
		require.Equal(t, 0.0, char.Rating())
	})

	t.Run("rating with factor", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, "СМ1-11Б")

		require.Zero(t, char.Rating())
		char.IncSkill(sm.Engineering, 3)
		char.IncSkill(sm.Engineering, 4)
		char.IncSkill(sm.Creative, 3)
		char.IncSkill(sm.Creative, 3)
		char.IncSkill(sm.Sportive, 4)

		expected := (3 + 4) * (1 + 0.3 + 0.3 + 0.4)
		require.InDelta(t, expected, char.Rating(), 1e-3)
	})
}

func TestCharacter_Booking(t *testing.T) {
	t.Run("should book Location", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, randomGroupName)
		require.NoError(t, char.Start())

		bookedTime := sm.MustNewBookedTime(
			char.Username,
			randomActivityUUID,
			timeWithMinutes(0),
			true,
		)
		err := char.AddBooking(bookedTime)
		require.NoError(t, err)
		require.True(t, char.HasBooking())
	})

	t.Run("should return error if character already has booking", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, randomGroupName)
		require.NoError(t, char.Start())

		bookedTime := sm.MustNewBookedTime(
			char.Username,
			randomActivityUUID,
			timeWithMinutes(0),
			true,
		)
		err := char.AddBooking(bookedTime)
		require.NoError(t, err)

		bookedTime = sm.MustNewBookedTime(
			char.Username,
			randomActivityUUID,
			timeWithMinutes(30),
			true,
		)
		err = char.AddBooking(bookedTime)
		require.ErrorIs(t, err, sm.ErrCharacterAlreadyHasBooking)
	})

	t.Run("should return error if booking is too late", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, randomGroupName)
		require.NoError(t, char.Start())

		start := time.Now().Add(-16 * time.Minute).Round(30 * time.Minute)
		bookedTime := sm.MustNewBookedTime(char.Username, randomActivityUUID, start, true)

		err := char.AddBooking(bookedTime)
		require.ErrorIs(t, err, sm.ErrCharacterBookingIsTooClose)
	})

	t.Run("should return error if character has ended Instruction", func(t *testing.T) {
		char := sm.MustNewCharacter(randomUsername, randomGroupName)
		require.NoError(t, char.Start())

		start := time.Now().Add(sm.MaxDurationInstruction).Round(20 * time.Minute).Add(20 * time.Minute)
		bookedTime := sm.MustNewBookedTime(char.Username, randomActivityUUID, start, true)

		err := char.AddBooking(bookedTime)
		require.ErrorIs(t, err, sm.ErrCharacterBookingAfterFinish)
	})
}

func TestCharacter_Start(t *testing.T) {
	char := sm.MustNewCharacter(randomUsername, randomGroupName)

	require.NoError(t, char.Start())
	require.True(t, char.IsStarted())
	require.True(t, char.IsProcessing())
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
