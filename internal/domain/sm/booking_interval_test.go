package sm_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"sm-instruction/internal/common/commonerrs"
	"sm-instruction/internal/domain/sm"
)

const randomUserID = 42

func TestBookingIntervalFactory_NewBookingInterval(t *testing.T) {
	duration := 10 * time.Minute
	cfg := sm.BookingIntervalFactoryConfig{
		IntervalDuration:             duration,
		MinimalDurationBeforeBooking: 5 * time.Minute,
	}
	factory := sm.MustNewBookingIntervalFactory(cfg)

	t.Run("should return a booking interval", func(t *testing.T) {
		start := timeWithMinutes(20)
		i, err := factory.NewBookingInterval(start, randomUserID)
		require.NoError(t, err)
		require.Equal(t, i.From.Hour(), start.Local().Hour())
		require.Equal(t, i.To.Sub(i.From), duration)
	})

	t.Run("should return error if booking interval is invalid", func(t *testing.T) {
		start := timeWithMinutes(15)
		_, err := factory.NewBookingInterval(start, randomUserID)
		require.ErrorAs(t, err, &commonerrs.InvalidInputError{})
	})
}

func TestBookingInterval_IsIntersects(t *testing.T) {
	cfg := sm.BookingIntervalFactoryConfig{
		IntervalDuration:             10 * time.Minute,
		MinimalDurationBeforeBooking: 1 * time.Second,
	}
	factory1 := sm.MustNewBookingIntervalFactory(cfg)

	cfg = sm.BookingIntervalFactoryConfig{
		IntervalDuration:             20 * time.Minute,
		MinimalDurationBeforeBooking: 1 * time.Second,
	}
	factory2 := sm.MustNewBookingIntervalFactory(cfg)

	t.Run("should return true if booking interval is intersects", func(t *testing.T) {
		/*
		             13:10      13:20
		               a          b
		    |----------|----------|----------|
		    c                     d
		  13:00                 13:20
		*/
		a := timeWithMinutes(10)
		c := timeWithMinutes(0)
		ab := factory1.MustNewBookingInterval(a, randomUserID)
		cd := factory2.MustNewBookingInterval(c, randomUserID)
		require.True(t, ab.IsIntersects(cd))
	})

	t.Run("should return true if booking intervals is equal", func(t *testing.T) {
		/*
		  13:10      13:20
		    a          b
		    |----------|
		    c          d
		  13:10      13:20
		*/
		a := timeWithMinutes(10)
		c := timeWithMinutes(10)
		ab := factory1.MustNewBookingInterval(a, randomUserID)
		cd := factory1.MustNewBookingInterval(c, randomUserID)
		require.True(t, ab.IsIntersects(cd))
	})

	t.Run("should return false if booking intervals are not intersects", func(t *testing.T) {
		/*
		  13:10      13:20
		    a          b
		    |----------|----------|
		               c          d
		             13:20      13:30
		*/
		a := timeWithMinutes(10)
		c := timeWithMinutes(20)
		ab := factory1.MustNewBookingInterval(a, randomUserID)
		cd := factory1.MustNewBookingInterval(c, randomUserID)
		require.False(t, ab.IsIntersects(cd))
	})

	t.Run("should return false if booking intervals in different hours", func(t *testing.T) {
		/*
		  13:10      13:20
		    a          b
		    |----------|----------|--- ... ---|----------|----------|
		                                                 c          d
		                                               14:20      14:30
		*/
		a := timeWithMinutes(10)
		c := timeWithMinutes(20).Add(time.Hour)
		require.Equal(t, time.Hour+10*time.Minute, c.Sub(a))
		ab := factory1.MustNewBookingInterval(a, randomUserID)
		cd := factory1.MustNewBookingInterval(c, randomUserID)
		require.False(t, ab.IsIntersects(cd))
		require.False(t, cd.IsIntersects(ab))
	})
}

func TestBookingIntervalFactory_AvailableIntervals(t *testing.T) {
	cfg := sm.BookingIntervalFactoryConfig{
		IntervalDuration:             20 * time.Minute,
		MinimalDurationBeforeBooking: 1 * time.Second,
	}
	factory := sm.MustNewBookingIntervalFactory(cfg)

	t.Run("should return available intervals", func(t *testing.T) {
		loc := sm.MustNewLocation("1234", "Test", []sm.SkillType{sm.Researching, sm.Creative})

		bookedTime := timeWithMinutes(20)
		booked := factory.MustNewBookingInterval(bookedTime, randomUserID)
		require.NoError(t, loc.AddBooking(booked))

		user := sm.MustNewUser(randomUserID, "test")
		char := sm.MustNewCharacter(user, randomGroupName)
		require.NoError(t, char.Start())

		available, err := factory.AvailableIntervals(char, loc)
		for _, av := range available {
			t.Log(av)
		}
		require.NoError(t, err)
		require.Truef(t, len(available) == 3*4-1 || len(available) == 3*4-2,
			"expected len 14 or 13, got %d", len(available),
		)
	})
}
