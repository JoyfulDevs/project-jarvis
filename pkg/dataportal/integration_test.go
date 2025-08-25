package dataportal_test

import (
	"context"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/devafterdark/project-jarvis/pkg/dataportal"
)

func TestListHolidays(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "on" {
		t.Skip("skipping integration test")
	}

	key, ok := os.LookupEnv("DATA_PORTAL_AUTH_KEY")
	if !ok {
		t.Fatal("DATA_PORTAL_AUTH_KEY is not set")
	}
	client := dataportal.NewClient(key)
	testCases := []struct {
		desc     string
		request  *dataportal.ListHolidaysRequest
		expected []string
	}{
		{
			desc: "2025년 1월",
			request: &dataportal.ListHolidaysRequest{
				Year:  2025,
				Month: 1,
			},
			expected: []string{
				"20250101", // 1월1일
				"20250127", // 임시공휴일
				"20250128", // 설날
				"20250129", // 설날
				"20250130", // 설날
			},
		},
		{
			desc: "2025년 2월",
			request: &dataportal.ListHolidaysRequest{
				Year:  2025,
				Month: 2,
			},
			expected: nil,
		},
		{
			desc: "2025년 3월",
			request: &dataportal.ListHolidaysRequest{
				Year:  2025,
				Month: 3,
			},
			expected: []string{
				"20250301", // 삼일절
				"20250303", // 대체공휴일
			},
		},
		{
			desc: "2025년 5월",
			request: &dataportal.ListHolidaysRequest{
				Year:  2025,
				Month: 5,
			},
			expected: []string{
				"20250505", // 어린이날
				"20250505", // 부처님오신날
				"20250506", // 대체공휴일
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			resp, err := client.ListHolidays(context.Background(), tc.request)
			if err != nil {
				t.Fatalf("failed to list holidays err=%v", err)
			}
			holidays := make([]string, 0, len(resp.Body.Data.Items))
			for _, item := range resp.Body.Data.Items {
				if item.IsHoliday != "Y" {
					continue
				}
				holidays = append(holidays, item.Date)
			}
			if len(holidays) != len(tc.expected) {
				t.Fatalf("invalid number of holidays value=%d, expected=%d", len(holidays), len(tc.expected))
			}
			for _, holiday := range holidays {
				if !slices.Contains(tc.expected, holiday) {
					t.Fatalf("unexpected holiday date=%s", holiday)
				}
			}
			for _, day := range tc.expected {
				contains := slices.ContainsFunc(holidays, func(e string) bool {
					return e == day
				})
				if !contains {
					t.Fatalf("holiday not found day=%s", day)
				}
			}
		})
	}
}

func TestGetUltraShortTermForecast(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "on" {
		t.Skip("skipping integration test")
	}

	key, ok := os.LookupEnv("DATA_PORTAL_AUTH_KEY")
	if !ok {
		t.Fatal("DATA_PORTAL_AUTH_KEY is not set")
	}

	client := dataportal.NewClient(key)

	// 현재 시간 기준으로 baseDate와 baseTime 설정
	// 초단기예보는 매시 30분에 발표되므로, 현재 시간보다 이전 시간으로 설정
	now := time.Now()
	baseDate := now.Format("20060102")

	// 현재 분이 30분 미만이면 이전 시간의 30분, 30분 이상이면 현재 시간의 30분
	hour := now.Hour()
	if now.Minute() < 30 {
		hour = hour - 1
		if hour < 0 {
			hour = 23
			now = now.AddDate(0, 0, -1)
			baseDate = now.Format("20060102")
		}
	}
	baseTime := fmt.Sprintf("%02d30", hour)

	testCases := []struct {
		desc    string
		request *dataportal.UltraShortTermForecastRequest
	}{
		{
			desc: "서울시청 좌표 초단기예보 조회",
			request: &dataportal.UltraShortTermForecastRequest{
				Page:     1,
				Count:    10,
				BaseDate: baseDate,
				BaseTime: baseTime,
				NX:       60,
				NY:       127,
			},
		},
		{
			desc: "부산시청 좌표 초단기예보 조회",
			request: &dataportal.UltraShortTermForecastRequest{
				Page:     1,
				Count:    20,
				BaseDate: baseDate,
				BaseTime: baseTime,
				NX:       98,
				NY:       76,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			resp, err := client.GetUltraShortTermForecast(context.Background(), tc.request)
			if err != nil {
				t.Fatalf("failed to get ultra short term forecast err=%v", err)
			}

			// 응답 코드가 성공인지 확인
			if resp.Header.Code != "00" {
				t.Fatalf("invalid response code=%s, message=%s", resp.Header.Code, resp.Header.Message)
			}

			// 데이터가 있는지 확인
			if len(resp.Body.Data.Items) == 0 {
				t.Fatal("empty response data")
			}

			// 좌표가 요청한 것과 동일한지 확인
			for _, item := range resp.Body.Data.Items {
				if item.NX != tc.request.NX || item.NY != tc.request.NY {
					t.Fatalf("invalid nx=%d, ny=%d", item.NX, item.NY)
				}
			}

			// 데이터 개수 확인 (요청한 Count보다 작거나 같아야 함)
			if len(resp.Body.Data.Items) > tc.request.Count {
				t.Fatalf("invalid response data count=%d, expected<=%d", len(resp.Body.Data.Items), tc.request.Count)
			}
		})
	}
}
