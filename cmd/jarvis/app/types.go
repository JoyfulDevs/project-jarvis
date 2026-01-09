package app

type Command = string

const (
	CommandEmpty           Command = ""
	CommandManual          Command = "기능"
	CommandHolidayCalendar Command = "공휴일"
	CommandForecast        Command = "날씨"
)

type ButtonAction = string

const (
	ButtonActionDone            ButtonAction = "done"
	ButtonActionManual          ButtonAction = "manual"
	ButtonActionHolidayCalendar ButtonAction = "holiday"
	ButtonActionForecast        ButtonAction = "forecast"
)
