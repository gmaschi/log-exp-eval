package str_test

import (
	"github.com/gmaschi/log-exp-eval/pkg/tools/str"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHasLowers(t *testing.T) {
	testCases := []struct {
		name              string
		str               string
		expectedLowers    []string
		expectedHasLowers bool
	}{
		{
			name:              "Has multiple lowers",
			str:               "-A e E j k LO)",
			expectedLowers:    []string{"e", "j", "k"},
			expectedHasLowers: true,
		},
		{
			name:              "Empty string",
			str:               "",
			expectedLowers:    []string{},
			expectedHasLowers: false,
		},
		{
			name:              "Does not have lowers",
			str:               "-A FE & * U A)",
			expectedLowers:    []string{},
			expectedHasLowers: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lowers, hasLowers := str.HasLowers(tc.str)
			require.Equal(t, tc.expectedHasLowers, hasLowers)
			require.Equal(t, tc.expectedLowers, lowers)
		})
	}
}
