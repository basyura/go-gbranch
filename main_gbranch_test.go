package main

import (
	"strings"
	"testing"
)

func TestAdjustSpace(t *testing.T) {
	type testCase struct {
		name       string
		s          string
		maxLen     int
		want       string
		wantWidth  int
		checkExact bool
	}

	tests := []testCase{
		{
			name:       "ascii padding",
			s:          "abc",
			maxLen:     5,
			want:       "abc  ",
			wantWidth:  5,
			checkExact: true,
		},
		{
			name:       "non-ascii chop avoids panic",
			s:          "あい",
			maxLen:     3,
			want:       "あ",
			wantWidth:  2,
			checkExact: true,
		},
		{
			name:       "mixed width exact fit",
			s:          "aあ",
			maxLen:     3,
			want:       "aあ",
			wantWidth:  3,
			checkExact: true,
		},
		{
			name:       "maxLen zero",
			s:          "abc",
			maxLen:     0,
			want:       "",
			wantWidth:  0,
			checkExact: true,
		},
		{
			name:      "maxLen capped by commitMsgLength",
			s:         strings.Repeat("a", 80),
			maxLen:    100,
			wantWidth: commitMsgLength,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := adjustSpace(tc.s, tc.maxLen)
			if tc.checkExact && got != tc.want {
				t.Fatalf("adjustSpace result mismatch: got %q, want %q", got, tc.want)
			}
			if tc.wantWidth >= 0 {
				if gotWidth := strLen(got); gotWidth != tc.wantWidth {
					t.Fatalf("adjustSpace width mismatch: got %d, want %d", gotWidth, tc.wantWidth)
				}
			}
		})
	}
}
