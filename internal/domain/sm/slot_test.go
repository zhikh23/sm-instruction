package sm_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

func TestSlot_Take(t *testing.T) {
	slot, err := sm.NewSlot(todayTime(11, 0), todayTime(11, 20))
	require.NoError(t, err)

	require.True(t, slot.IsAvailable())

	err = slot.Take("СМ1-11Б")
	require.NoError(t, err)

	require.False(t, slot.IsAvailable())

	err = slot.Take("СМ2-12")
	require.ErrorIs(t, err, sm.ErrSlotHasAlreadyTaken)

	err = slot.Free()
	require.NoError(t, err)

	require.True(t, slot.IsAvailable())

	err = slot.Free()
	require.ErrorIs(t, err, sm.ErrSlotHasNotTaken)
}

func todayTime(hours int, minutes int) time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hours, minutes, 0, 0, time.Local)
}
