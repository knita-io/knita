package file

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsGlob(t *testing.T) {
	scenarios := []struct {
		input string
		glob  bool
		err   error
	}{
		{"?", true, nil},
		{"fooba?", true, nil},
		{"foobar", false, nil},
		{"*", true, nil},
		{"**", true, nil},
		{"foobar/*", true, nil},
		{"foobar/**", true, nil},
		{"foobar/{a}", false, nil},
		{"foobar/{a,b}", true, nil},
		{"foobar/{a\\dd}", false, nil},
		{"foobar/{a\\dd ", false, nil},
		{"foobar/{a\\dd ?", true, nil},
		{"foobar/{a\\dd,}", true, nil},
		{"foobar/[ ", false, nil},
		{"foobar/[a]", true, nil},
		{"foobar/[a-z]", true, nil},
		{"foobar/[abc]", true, nil},
		{"foobar/[^abc]", true, nil},
		{"foobar/[!abc]", true, nil},
		{"foobar/[ ]", false, nil},
		{"foobar/[ d]", false, nil},
		{"foobar/[d ]", true, nil},
	}
	for _, scenario := range scenarios {
		glob, err := isGlob(scenario.input)
		if scenario.err != nil {
			require.Equal(t, scenario.err, err)
		} else {
			require.NoError(t, err)
		}
		require.Equal(t, scenario.glob, glob)
	}
}
