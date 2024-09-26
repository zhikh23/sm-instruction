package sm_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"sm-instruction/internal/domain/sm"
)

func TestSlot_Take(t *testing.T) {
	slot, err := sm.NewSlot(time.Now(), time.Now().Add(20*time.Minute))
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
