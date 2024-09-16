package sm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"sm-instruction/internal/domain/sm"
)

func TestLocation_AddBooking(t *testing.T) {
	t.Run("should add booking", func(t *testing.T) {
		loc := sm.MustNewLocation("Desc", "ИЦАР")

		bookedTime := sm.MustNewBookedTime(
			randomUsername,
			randomActivityUUID,
			timeWithMinutes(0),
			true,
		)
		err := loc.AddBooking(bookedTime)
		require.NoError(t, err)
		require.True(t, loc.IsBooked(bookedTime.Start))
		require.False(t, loc.IsBooked(bookedTime.Finish))
		require.ErrorIs(t, loc.CanBook(bookedTime.Start), sm.ErrLocationAlreadyBooked)
		require.NoError(t, loc.CanBook(bookedTime.Finish))
	})

	t.Run("should return error if already booked", func(t *testing.T) {
		loc := sm.MustNewLocation("Desc", "509м")

		bookedTime := sm.MustNewBookedTime(
			"1",
			randomActivityUUID,
			timeWithMinutes(0),
			true,
		)
		err := loc.AddBooking(bookedTime)
		require.NoError(t, err)

		bookedTime = sm.MustNewBookedTime(
			"2",
			randomActivityUUID,
			timeWithMinutes(0),
			true,
		)
		err = loc.AddBooking(bookedTime)
		require.ErrorIs(t, err, sm.ErrLocationAlreadyBooked)
	})
}

func TestLocation_RemoveBooking(t *testing.T) {
	t.Run("should remove booking", func(t *testing.T) {
		loc := sm.MustNewLocation("Desc", "509м")

		bookedTime := sm.MustNewBookedTime(
			"1",
			randomActivityUUID,
			timeWithMinutes(0),
			true,
		)
		err := loc.AddBooking(bookedTime)
		require.NoError(t, err)

		err = loc.RemoveBooking(bookedTime)
		require.NoError(t, err)

		require.False(t, loc.IsBooked(bookedTime.Start))
		require.NoError(t, loc.CanBook(bookedTime.Start))
	})

	t.Run("should return error if booking is not exists", func(t *testing.T) {
		loc := sm.MustNewLocation("Desc", "509м")

		bookedTime := sm.MustNewBookedTime(
			"1",
			randomActivityUUID,
			timeWithMinutes(0),
			true,
		)
		err := loc.RemoveBooking(bookedTime)
		require.Error(t, err, sm.ErrLocationHasNotBooked)
	})
}
