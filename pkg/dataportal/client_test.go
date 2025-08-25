package dataportal_test

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/devafterdark/project-jarvis/pkg/dataportal"
)

func TestUltraShortTermForecastResponse(t *testing.T) {
	mockJSON := `{
		"response": {
			"header": {
				"resultCode": "00",
				"resultMsg": "NORMAL_SERVICE"
			},
			"body": {
				"dataType": "JSON",
				"items": {
					"item": [
						{
							"baseDate": "20250628",
							"baseTime": "1430",
							"category": "T1H",
							"fcstDate": "20250628",
							"fcstTime": "1500",
							"fcstValue": "25",
							"nx": 60,
							"ny": 127
						},
						{
							"baseDate": "20250628",
							"baseTime": "1430",
							"category": "SKY",
							"fcstDate": "20250628",
							"fcstTime": "1500",
							"fcstValue": "1",
							"nx": 60,
							"ny": 127
						}
					]
				},
				"pageNo": 1,
				"numOfRows": 10,
				"totalCount": 2
			}
		}
	}`

	resp := &dataportal.UltraShortTermForecastResponse{}
	err := json.Unmarshal([]byte(mockJSON), resp)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// 헤더 검증
	if resp.Header.Code != "00" {
		t.Errorf("Header.Code: got %s, want 00", resp.Header.Code)
	}
	if resp.Header.Message != "NORMAL_SERVICE" {
		t.Errorf("Header.Message: got %s, want NORMAL_SERVICE", resp.Header.Message)
	}

	// 바디 검증
	if resp.Body.Page != 1 {
		t.Errorf("Body.Page: got %d, want 1", resp.Body.Page)
	}
	if resp.Body.Count != 10 {
		t.Errorf("Body.Count: got %d, want 10", resp.Body.Count)
	}
	if resp.Body.Total != 2 {
		t.Errorf("Body.Total: got %d, want 2", resp.Body.Total)
	}

	// 아이템 검증
	if len(resp.Body.Data.Items) != 2 {
		t.Fatalf("Items 길이: got %d, want 2", len(resp.Body.Data.Items))
	}

	// 첫 번째 아이템 (기온)
	item1 := resp.Body.Data.Items[0]
	if item1.Category != dataportal.CategoryTemperature {
		t.Errorf("Item1.Category: got %s, want %s", item1.Category, dataportal.CategoryTemperature)
	}
	if item1.ForecastValue != "25" {
		t.Errorf("Item1.ForecastValue: got %s, want 25", item1.ForecastValue)
	}
	if item1.NX != 60 {
		t.Errorf("Item1.NX: got %d, want 60", item1.NX)
	}
	if item1.NY != 127 {
		t.Errorf("Item1.NY: got %d, want 127", item1.NY)
	}

	// 두 번째 아이템 (하늘상태)
	item2 := resp.Body.Data.Items[1]
	if item2.Category != dataportal.CategorySky {
		t.Errorf("Item2.Category: got %s, want %s", item2.Category, dataportal.CategorySky)
	}
	if item2.ForecastValue != string(dataportal.SkyCodeClear) {
		t.Errorf("Item2.ForecastValue: got %s, want %s", item2.ForecastValue, dataportal.SkyCodeClear)
	}
}

func TestListHolidaysResponse(t *testing.T) {
	mockXML := `
	<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
	<response>
		<header>
			<resultCode>00</resultCode>
			<resultMsg>NORMAL SERVICE.</resultMsg>
		</header>
		<body>
			<items>
				<item>
					<dateKind>01</dateKind>
					<dateName>어린이날</dateName>
					<isHoliday>Y</isHoliday>
					<locdate>20250505</locdate>
					<seq>1</seq>
				</item>
				<item>
					<dateKind>01</dateKind>
					<dateName>부처님오신날</dateName>
					<isHoliday>Y</isHoliday>
					<locdate>20250505</locdate>
					<seq>2</seq>
				</item>
				<item>
					<dateKind>01</dateKind>
					<dateName>대체공휴일</dateName>
					<isHoliday>Y</isHoliday>
					<locdate>20250506</locdate>
					<seq>1</seq>
				</item>
			</items>
			<numOfRows>10</numOfRows>
			<pageNo>1</pageNo>
			<totalCount>3</totalCount>
		</body>
	</response>`

	resp := &dataportal.ListHolidaysResponse{}
	err := xml.Unmarshal([]byte(mockXML), resp)
	if err != nil {
		t.Fatalf("failed to unmarshal XML: %v", err)
	}

	// 헤더 검증
	if resp.Header.Code != "00" {
		t.Errorf("Header.Code: got %s, want 00", resp.Header.Code)
	}
	if resp.Header.Message != "NORMAL SERVICE." {
		t.Errorf("Header.Message: got %s, want NORMAL SERVICE.", resp.Header.Message)
	}

	// 바디 검증
	if resp.Body.Page != 1 {
		t.Errorf("Body.Page: got %d, want 1", resp.Body.Page)
	}
	if resp.Body.Count != 10 {
		t.Errorf("Body.Count: got %d, want 10", resp.Body.Count)
	}
	if resp.Body.Total != 3 {
		t.Errorf("Body.Total: got %d, want 3", resp.Body.Total)
	}

	// 아이템 검증
	if len(resp.Body.Data.Items) != 3 {
		t.Fatalf("Items 길이: got %d, want 3", len(resp.Body.Data.Items))
	}

	// 첫 번째 아이템 (어린이날)
	item1 := resp.Body.Data.Items[0]
	if item1.Name != "어린이날" {
		t.Errorf("Item1.Name: got %s, want 어린이날", item1.Name)
	}
	if item1.Date != "20250505" {
		t.Errorf("Item1.Date: got %s, want 20250505", item1.Date)
	}
	if item1.IsHoliday != "Y" {
		t.Errorf("Item1.IsHoliday: got %s, want Y", item1.IsHoliday)
	}

	// 두 번째 아이템 (부처님오신날)
	item2 := resp.Body.Data.Items[1]
	if item2.Name != "부처님오신날" {
		t.Errorf("Item2.Name: got %s, want 부처님오신날", item2.Name)
	}
	if item2.Date != "20250505" {
		t.Errorf("Item2.Date: got %s, want 20250505", item2.Date)
	}
	if item2.IsHoliday != "Y" {
		t.Errorf("Item2.IsHoliday: got %s, want Y", item2.IsHoliday)
	}

	// 세 번째 아이템 (대체공휴일)
	item3 := resp.Body.Data.Items[2]
	if item3.Name != "대체공휴일" {
		t.Errorf("Item3.Name: got %s, want 대체공휴일", item3.Name)
	}
	if item3.Date != "20250506" {
		t.Errorf("Item3.Date: got %s, want 20250506", item3.Date)
	}
	if item3.IsHoliday != "Y" {
		t.Errorf("Item3.IsHoliday: got %s, want Y", item3.IsHoliday)
	}
}
