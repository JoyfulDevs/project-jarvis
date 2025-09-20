package kst_test

import (
	"testing"
	"time"

	"github.com/joyfuldevs/project-jarvis/pkg/kst"
)

func TestKoreaTime(t *testing.T) {
	testCases := []struct {
		desc  string
		input time.Time
	}{
		{
			desc:  "UTC",
			input: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "Local",
			input: time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
		},
		{
			desc:  "KST",
			input: kst.KST(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			k := kst.KST(tc.input)
			if _, offset := kst.KST(tc.input).Zone(); offset != 9*60*60 {
				t.Fatalf("expected KST offset to be 9 hours, got %d seconds", offset)
			}
			if k.UTC() != tc.input.UTC() {
				t.Fatalf("expected KST time to be %s, got %s", tc.input.UTC(), k.UTC())
			}
		})
	}
}

func TestWeekday(t *testing.T) {
	testCases := []struct {
		desc     string
		input    time.Time
		expected string
	}{
		{
			desc:     "25년 1월 1일 (수)",
			input:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "수",
		},
		{
			desc:     "25년 1월 2일 (목)",
			input:    time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			expected: "목",
		},
		{
			desc:     "25년 1월 3일 (금)",
			input:    time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
			expected: "금",
		},
		{
			desc:     "25년 1월 4일 (토)",
			input:    time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
			expected: "토",
		},
		{
			desc:     "25년 1월 5일 (일)",
			input:    time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
			expected: "일",
		},
		{
			desc:     "25년 1월 6일 (월)",
			input:    time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),
			expected: "월",
		},
		{
			desc:     "25년 1월 7일 (화)",
			input:    time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC),
			expected: "화",
		},
		{
			desc:     "25년 1월 8일 (수)",
			input:    time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC),
			expected: "수",
		},
		{
			desc:     "25년 1월 9일 (목)",
			input:    time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC),
			expected: "목",
		},
		{
			desc:     "25년 1월 10일 (금)",
			input:    time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
			expected: "금",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if wd := kst.Weekday(tc.input); wd != tc.expected {
				t.Fatalf("expected weekday to be %s, got %s", tc.expected, wd)
			}
		})
	}
}

func TestLastWeekday(t *testing.T) {
	testCases := []struct {
		desc  string
		input time.Time
	}{
		{
			desc:  "25년 1월 1일 (수)",
			input: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 2일 (목)",
			input: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 3일 (금)",
			input: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 4일 (토)",
			input: time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 5일 (일)",
			input: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 6일 (월)",
			input: time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 7일 (화)",
			input: time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 8일 (수)",
			input: time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 9일 (목)",
			input: time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 10일 (금)",
			input: time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			for _, wd := range []time.Weekday{
				time.Monday,
				time.Tuesday,
				time.Wednesday,
				time.Thursday,
				time.Friday,
				time.Saturday,
				time.Sunday,
			} {
				result := kst.LastWeekday(tc.input, wd)
				if weekday := result.Weekday(); weekday != wd {
					t.Fatalf("expected last weekday to be %s, got %s", wd, weekday)
				}
				if result.After(tc.input) {
					t.Fatalf("expected last weekday to be before input time, got %s before %s", tc.input, result)
				}
				if tc.input.After(result.AddDate(0, 0, 7)) {
					t.Fatalf("expected last weekday +7 days to be after input time, got %s after %s", result.AddDate(0, 0, 7), tc.input)
				}
			}
		})
	}
}

func TestNextWeekday(t *testing.T) {
	testCases := []struct {
		desc  string
		input time.Time
	}{
		{
			desc:  "25년 1월 1일 (수)",
			input: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 2일 (목)",
			input: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 3일 (금)",
			input: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 4일 (토)",
			input: time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 5일 (일)",
			input: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 6일 (월)",
			input: time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 7일 (화)",
			input: time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 8일 (수)",
			input: time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 9일 (목)",
			input: time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 10일 (금)",
			input: time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			weekdays := []time.Weekday{
				time.Monday,
				time.Tuesday,
				time.Wednesday,
				time.Thursday,
				time.Friday,
				time.Saturday,
				time.Sunday,
			}
			for _, wd := range weekdays {
				result := kst.NextWeekday(tc.input, wd)
				if weekday := result.Weekday(); weekday != wd {
					t.Fatalf("expected next weekday to be %s, got %s", wd, weekday)
				}
				if tc.input.After(result) {
					t.Fatalf("expected next weekday to be after input time, got %s after %s", tc.input, result)
				}
				if result.AddDate(0, 0, -7).After(tc.input) {
					t.Fatalf("expected next weekday -7 days to be before input time, got %s before %s", tc.input, result.AddDate(0, 0, -7))
				}
			}
		})
	}
}
