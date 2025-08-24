package dataportal

type Response[T any] struct {
	Header Header   `json:"header" xml:"header"`
	Body   *Body[T] `json:"body,omitempty" xml:"body"`
}

type Header struct {
	// 결과 코드.
	Code string `json:"resultCode" xml:"resultCode"`
	// 결과 메시지.
	Message string `json:"resultMsg" xml:"resultMsg"`
}

type Body[T any] struct {
	// 페이지 번호.
	Page int `json:"pageNo" xml:"pageNo"`
	// 페이지 당 데이터 수.
	Count int `json:"numOfRows" xml:"numOfRows"`
	// 전체 데이터 수.
	Total int `json:"totalCount" xml:"totalCount"`
	Data  struct {
		Items []T `json:"item" xml:"item"`
	} `json:"items" xml:"items"`
}

// 공휴일 정보 조회
// https://www.data.go.kr/data/15012690/openapi.do

type ListHolidaysItem struct {
	// 날짜.
	// e.g. 20150301
	Date string `json:"locdate" xml:"locdate"`
	// 공공기관 휴일 여부.
	// e.g. Y
	IsHoliday string `json:"isHoliday" xml:"isHoliday"`
	// 명칭.
	// e.g. 삼일절
	Name string `json:"dateName" xml:"dateName"`
}

type ListHolidaysRequest struct {
	Year  int
	Month int
}

type ListHolidaysResponse struct {
	Response[ListHolidaysItem] `xml:"response"`
}

// 초단기 예보 조회
// https://www.data.go.kr/data/15084084/openapi.do

type UltraShortTermForecastCategory string

const (
	// 기온 ℃
	CategoryTemperature UltraShortTermForecastCategory = "T1H"
	// 1시간 강수량 범주 (1 mm)
	CategoryRainfall UltraShortTermForecastCategory = "RN1"
	// 하늘 상태 코드값
	CategorySky UltraShortTermForecastCategory = "SKY"
	// 동서바람성분 m/s
	CategoryEastWestWind UltraShortTermForecastCategory = "UUU"
	// 남북바람성분 m/s
	CategoryNorthSouthWind UltraShortTermForecastCategory = "VVV"
	// 습도 %
	CategoryHumidity UltraShortTermForecastCategory = "REH"
	// 강수 형태 코드값
	CategoryPrecipitation UltraShortTermForecastCategory = "PTY"
	// 낙뢰 kA(킬로암페어)
	CategoryLightning UltraShortTermForecastCategory = "LGT"
	// 풍향 deg
	CategoryWindDirection UltraShortTermForecastCategory = "VEC"
	// 풍속 m/s
	CategoryWindSpeed UltraShortTermForecastCategory = "WSD"
)

func (c UltraShortTermForecastCategory) String() string {
	switch c {
	case CategoryTemperature:
		return "기온"
	case CategoryRainfall:
		return "1시간 강수량"
	case CategorySky:
		return "하늘 상태"
	case CategoryEastWestWind:
		return "동서바람성분"
	case CategoryNorthSouthWind:
		return "남북바람성분"
	case CategoryHumidity:
		return "습도"
	case CategoryPrecipitation:
		return "강수 형태"
	case CategoryLightning:
		return "낙뢰"
	case CategoryWindDirection:
		return "풍향"
	case CategoryWindSpeed:
		return "풍속"
	default:
		return "?"
	}
}

type PrecipitationCode string

const (
	// 없음
	PrecipitationCodeNone PrecipitationCode = "0"
	// 비
	PrecipitationCodeRain PrecipitationCode = "1"
	// 비/눈
	PrecipitationCodeRainSnow PrecipitationCode = "2"
	// 눈
	PrecipitationCodeSnow PrecipitationCode = "3"
	// 소나기
	PrecipitationCodeShower PrecipitationCode = "4"
	// 빗방울
	PrecipitationCodeDrizzle PrecipitationCode = "5"
	// 진눈깨비
	PrecipitationCodeSleet PrecipitationCode = "6"
	// 눈날림
	PrecipitationCodeSnowGrain PrecipitationCode = "7"
)

func (p PrecipitationCode) String() string {
	switch p {
	case PrecipitationCodeNone:
		return "없음"
	case PrecipitationCodeRain:
		return "비"
	case PrecipitationCodeRainSnow:
		return "비/눈"
	case PrecipitationCodeSnow:
		return "눈"
	case PrecipitationCodeShower:
		return "소나기"
	case PrecipitationCodeDrizzle:
		return "빗방울"
	case PrecipitationCodeSleet:
		return "진눈깨비"
	case PrecipitationCodeSnowGrain:
		return "눈날림"
	default:
		return "?"
	}
}

type SkyCode string

const (
	// 맑음
	SkyCodeClear SkyCode = "1"
	// 구름 조금
	SkyCodeFewClouds SkyCode = "2"
	// 구름 많음
	SkyCodePartlyCloudy SkyCode = "3"
	// 흐림
	SkyCodeCloudy SkyCode = "4"
)

func (s SkyCode) String() string {
	switch s {
	case SkyCodeClear:
		return "맑음"
	case SkyCodeFewClouds:
		return "구름 조금"
	case SkyCodePartlyCloudy:
		return "구름 많음"
	case SkyCodeCloudy:
		return "흐림"
	default:
		return "?"
	}
}

type UltraShortTermForecastItem struct {
	// 발표 일자.
	// e.g. 20250101
	Date string `json:"baseDate" xml:"baseDate"`
	// 발표 시간.
	// e.g. 0930
	Time string `json:"baseTime" xml:"baseTime"`
	// 예보 구분.
	Category UltraShortTermForecastCategory `json:"category" xml:"category"`
	// 예측 일자.
	// e.g. 20250101
	ForecastDate string `json:"fcstDate" xml:"fcstDate"`
	// 예측 시간.
	// e.g. 1200
	ForecastTime string `json:"fcstTime" xml:"fcstTime"`
	// 예보 값.
	ForecastValue string `json:"fcstValue" xml:"fcstValue"`
	// 예보 지점 X 좌표.
	NX int `json:"nx" xml:"nx"`
	// 예보 지점 Y 좌표.
	NY int `json:"ny" xml:"ny"`
}

type UltraShortTermForecastRequest struct {
	// 페이지 번호.
	Page int
	// 페이지 당 데이터 수.
	Count int
	// 발표 일자.
	// e.g. 20250101
	BaseDate string
	// 발표 시간.
	// e.g. 0930
	BaseTime string
	// 예보 지점 X 좌표.
	NX int
	// 예보 지점 Y 좌표.
	NY int
}

type UltraShortTermForecastResponse struct {
	Response[UltraShortTermForecastItem] `json:"response"`
}
