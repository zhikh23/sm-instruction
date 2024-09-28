package sm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

func TestGrades(t *testing.T) {
	char := sm.MustNewCharacter("СМ1-11Б", "testname", []*sm.Slot{})
	require.Zero(t, char.Rating())

	var err error
	err = char.GiveGrade(sm.Engineering, 3, "ЦМР")
	require.NoError(t, err)
	require.Equal(t, 3.0, char.Rating())

	err = char.GiveGrade(sm.Researching, 2, "ЦМР")
	require.NoError(t, err)
	require.Equal(t, 5.0, char.Rating())

	err = char.GiveGrade(sm.Creative, 2, "ССФСМ")
	require.NoError(t, err)
	require.InDelta(t, 5.13889, char.Rating(), 1e-5)
}
