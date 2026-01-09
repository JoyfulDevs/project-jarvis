package app

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/joyfuldevs/project-jarvis/pkg/slack"
	"github.com/joyfuldevs/project-jarvis/pkg/slack/blockkit"
)

type CommandResponder struct {
	AppToken string
	BotToken string
	Payload  *slack.SlashCommandEventPayload
}

func (c *CommandResponder) RespondCommand() {
	switch c.Payload.Text {
	case CommandEmpty:
		c.RespondCommandEmpty()
	case CommandManual:
		c.RespondCommandManual()
	case CommandHolidayCalendar:
		c.RespondCommandHolidayCalendar()
	case CommandForecast:
		c.RespondCommandForecast()
	default:
		c.RespondCommandUndefined()
	}
}

func (c *CommandResponder) RespondCommandEmpty() {
	Respond(c.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeManualMessage(),
		ReplaceOriginal: true,
	})
}

func (c *CommandResponder) RespondCommandManual() {
	Respond(c.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeManualMessage(),
		ReplaceOriginal: true,
	})
}

func (c *CommandResponder) RespondCommandHolidayCalendar() {
	Respond(c.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeHolidayCalendarMessage(),
		ReplaceOriginal: true,
	})
}

func (c *CommandResponder) RespondCommandForecast() {
	Respond(c.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeForecastMessage(),
		ReplaceOriginal: true,
	})
}

func (c *CommandResponder) RespondCommandUndefined() {
	Respond(c.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeGuideMessage(),
		ReplaceOriginal: true,
	})
}

type ActionResponder struct {
	AppToken string
	BotToken string
	Payload  *slack.InteractiveEventPayload
}

func (a *ActionResponder) RespondActions() {
	for _, action := range a.Payload.Actions {
		a.RespondAction(action)
	}
}

func (a *ActionResponder) RespondAction(action slack.InteractiveAction) {
	switch action.Type {
	case blockkit.ElementTypeButton:
		// 버튼 액션 처리.
		a.RespondButtonAction(action)
	case blockkit.ElementTypeSelect:
		// 셀렉트 액션 처리.
	default:
		slog.Warn("undefined action type", slog.String("type", string(action.Type)))
	}
}

func (a *ActionResponder) RespondButtonAction(action slack.InteractiveAction) {
	switch action.ActionID {
	case ButtonActionDone:
		// 완료 버튼 클릭.
		a.RespondButtonActionDone()
	case ButtonActionManual:
		// 기능 안내 버튼 클릭.
		a.RespondButtonActionManual()
	case ButtonActionHolidayCalendar:
		// 공휴일 안내 버튼 클릭.
		a.RespondProgress()
		a.RespondButtonActionHolidayCalendar()
	case ButtonActionForecast:
		// 날씨 버튼 클릭.
		a.RespondProgress()
		a.RespondButtonActionForecast()
	}
}

func (a *ActionResponder) RespondProgress() {
	Respond(a.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeProgressMessage(),
		ReplaceOriginal: true,
	})
}

func (a *ActionResponder) RespondButtonActionDone() {
	Respond(a.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		DeleteOriginal: true,
	})
}

func (a *ActionResponder) RespondButtonActionManual() {
	Respond(a.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeManualMessage(),
		ReplaceOriginal: true,
	})
}

func (a *ActionResponder) RespondButtonActionHolidayCalendar() {
	Respond(a.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeHolidayCalendarMessage(),
		ReplaceOriginal: true,
	})
}

func (a *ActionResponder) RespondButtonActionForecast() {
	Respond(a.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeForecastMessage(),
		ReplaceOriginal: true,
	})
}

func Respond(url string, payload *slack.InteractiveResponsePayload) {
	bodyData, err := json.Marshal(payload)
	if err != nil {
		slog.Error("failed to respond", slog.Any("error", err))
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyData))
	if err != nil {
		slog.Error("failed to respond", slog.Any("error", err))
		return
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("failed to respond", slog.Any("error", err))
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return
	}

	if body, err := io.ReadAll(resp.Body); err == nil {
		msg := string(body)
		slog.Error("failed to respond", slog.Any("error", errors.New(msg)))
		if strings.Contains(msg, "invalid_blocks") {
			encoded := base64.StdEncoding.EncodeToString(bodyData)
			slog.Info("request body", slog.String("base64", encoded))
		}
	} else {
		slog.Error("failed to respond", slog.Any("error", errors.New(resp.Status)))
	}
}
