package kst

import "time"

var kst = time.FixedZone("KST", 60*60*9)

// 현재 한국 시간을 리턴한다.
func Now() time.Time {
	return time.Now().In(kst)
}

// 전달 받은 시간을 한국 시간으로 변환한다.
func KST(t time.Time) time.Time {
	return t.In(kst)
}

// 날짜에 따라 월,화,수,목,금,토,일 요일을 리턴한다.
func Weekday(t time.Time) string {
	switch t.Weekday() {
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

// 전달받은 날짜를 기준으로 마지막 월요일의 날짜를 리턴한다.
func LastMonday(t time.Time) time.Time {
	t = KST(t)
	weekday := t.Weekday()
	offset := int(time.Monday - weekday)
	if offset > 0 {
		offset -= 7
	}
	return time.Date(t.Year(), t.Month(), t.Day()+offset, 0, 0, 0, 0, t.Location())
}
