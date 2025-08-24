package app

type Command = string

const (
	CommandEmpty           Command = ""
	CommandManual          Command = "기능"
	CommandHolidayCalendar Command = "공휴일"
	CommandForecast        Command = "날씨"
	CommandScrumList       Command = "스크럼"
	CommandScrumSummary    Command = "스크럼요약"
	CommandConfig          Command = "설정"
)

type ButtonAction = string

const (
	ButtonActionDone            ButtonAction = "done"
	ButtonActionManual          ButtonAction = "manual"
	ButtonActionHolidayCalendar ButtonAction = "holiday"
	ButtonActionForecast        ButtonAction = "forecast"
	ButtonActionScrumList       ButtonAction = "scrum_list"
	ButtonActionScrumSummary    ButtonAction = "scrum_summary"
	ButtonActionConfig          ButtonAction = "config"
)

type ConfigAction = string

const (
	ConfigActionDailyScrumEnable   ConfigAction = "daily_scrum_enable"
	ConfigActionWeeklyReportEnable ConfigAction = "weekly_report_enable"
)
