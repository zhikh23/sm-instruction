package sm

import (
	"github.com/stretchr/testify/require"
	"sm-instruction/internal/common/commonerrs"
	"testing"
)

func TestNewSkillTypeFromString(t *testing.T) {
	t.Run("should return a new SkillType", func(t *testing.T) {
		type testCase struct {
			str string
			typ SkillType
		}

		cases := []testCase{
			{"engineering", Engineering},
			{"researching", Researching},
			{"social", Social},
			{"creative", Creative},
			{"sportive", Sportive},
		}

		for _, tc := range cases {
			st, err := NewSkillTypeFromString(tc.str)
			require.NoError(t, err)
			require.Equal(t, tc.typ, st)
			require.Equal(t, tc.str, st.String())
		}
	})

	t.Run("should return an error on invalid SkillType", func(t *testing.T) {
		st, err := NewSkillTypeFromString("")
		require.ErrorAs(t, err, &commonerrs.InvalidInputError{})
		require.True(t, st.IsZero())
	})
}
