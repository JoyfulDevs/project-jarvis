package kst

import "time"

// KST is UTC +09:00
const Offset = 9 * 60 * 60

var Zone = time.FixedZone("KST", Offset)

type weekday time.Weekday

func (w weekday) String() string {
	switch time.Weekday(w) {
	case time.Sunday:
		return "일"
	case time.Monday:
		return "월"
	case time.Tuesday:
		return "화"
	case time.Wednesday:
		return "수"
	case time.Thursday:
		return "목"
	case time.Friday:
		return "금"
	case time.Saturday:
		return "토"
	default:
		return ""
	}
}

// 한국 표준시(KST)로 설정된 현재 시간을 반환한다.
func Now() time.Time {
	return time.Now().In(Zone)
}

// 주어진 시간을 한국 표준시(KST)로 변환한다.
func KST(t time.Time) time.Time {
	return t.In(Zone)
}

// 주어진 시간의 한국어 요일을 반환한다.
func Weekday(t time.Time) string {
	return weekday(t.Weekday()).String()
}

// 주어진 날짜를 기준으로 가장 최근의 특정 요일의 시간을 반환한다.
func LastWeekday(t time.Time, weekday time.Weekday) time.Time {
	t = KST(t)
	year, month, day := t.Date()
	offset := int(weekday - t.Weekday())
	if offset > 0 {
		offset -= 7
	}
	return time.Date(year, month, day+offset, 0, 0, 0, 0, t.Location())
}

// 주어진 날짜를 기준으로 가장 가까운 특정 요일의 시간을 반환한다.
func NextWeekday(t time.Time, weekday time.Weekday) time.Time {
	t = KST(t)
	year, month, day := t.Date()
	offset := int(weekday - t.Weekday())
	if offset <= 0 {
		offset += 7
	}
	return time.Date(year, month, day+offset, 0, 0, 0, 0, t.Location())
}
