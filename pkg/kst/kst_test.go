package kst_test

import (
	"testing"
	"time"

	"github.com/genians/endpoint-lab-slack-bot/pkg/kst"
)

func TestKoreaTime(t *testing.T) {
	testCases := []struct {
		desc     string
		input    time.Time
		expected time.Time
	}{
		{
			desc:     "UTC 시간을 한국 시간으로 변환",
			input:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Hour * 9),
		},
		{
			desc:     "한국 시간을 다시 한국 시간으로 변환",
			input:    kst.KST(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
			expected: kst.KST(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			kst := kst.KST(tc.input)
			if kst.Year() != tc.expected.Year() {
				t.Errorf("invalid year value=%d, expected=%d", kst.Year(), tc.expected.Year())
			}
			if kst.Month() != tc.expected.Month() {
				t.Errorf("invalid month value=%d, expected=%d", kst.Month(), tc.expected.Month())
			}
			if kst.Day() != tc.expected.Day() {
				t.Errorf("invalid day value=%d, expected=%d", kst.Day(), tc.expected.Day())
			}
			if kst.Hour() != tc.expected.Hour() {
				t.Errorf("invalid hour value=%d, expected=%d", kst.Hour(), tc.expected.Hour())
			}
			if kst.Minute() != tc.expected.Minute() {
				t.Errorf("invalid minute value=%d, expected=%d", kst.Minute(), tc.expected.Minute())
			}
			if kst.Second() != tc.expected.Second() {
				t.Errorf("invalid second value=%d, expected=%d", kst.Second(), tc.expected.Second())
			}
		})
	}
}

func TestLastMonday(t *testing.T) {
	testCases := []struct {
		desc  string
		input time.Time
	}{
		{
			desc:  "25년 1월 1일",
			input: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 2일",
			input: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 3일",
			input: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 4일",
			input: time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 5일",
			input: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 6일",
			input: time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 7일",
			input: time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 8일",
			input: time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 9일",
			input: time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:  "25년 1월 10일",
			input: time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := kst.LastMonday(tc.input)
			if result.Weekday() != time.Monday {
				t.Fatalf("is not monday time=%s", result)
			}
		})
	}
}
