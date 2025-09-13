package app

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/joyfuldevs/project-jarvis/pkg/kst"
	"github.com/joyfuldevs/project-jarvis/pkg/slack"
	"github.com/joyfuldevs/project-jarvis/pkg/slack/blockkit"
	channelconfig "github.com/joyfuldevs/project-jarvis/service/channelconfig/client"
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
	case CommandScrumList:
		c.RespondCommandScrumList()
	case CommandScrumSummary:
		c.RespondCommandScrumSummary()
	case CommandConfig:
		c.RespondCommandConfig()
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

func (c *CommandResponder) RespondCommandScrumList() {
	messages := listScrumMessages(c.AppToken, c.BotToken, c.Payload.ChannelID, c.Payload.UserID)
	Respond(c.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeScrumListMessage(messages),
		ReplaceOriginal: true,
	})
}

func (c *CommandResponder) RespondCommandScrumSummary() {
	messages := listScrumMessages(c.AppToken, c.BotToken, c.Payload.ChannelID, c.Payload.UserID)
	Respond(c.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeScrumSummaryMessage(messages),
		ReplaceOriginal: true,
	})
}

func (c *CommandResponder) RespondCommandConfig() {
	Respond(c.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeConfigMessage(c.Payload.ChannelID),
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
		a.RespondSelectAction(action)
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
	case ButtonActionScrumList:
		// 스크럼 목록 버튼 클릭.
		a.RespondProgress()
		a.RespondButtonActionScrumList()
	case ButtonActionScrumSummary:
		// 스크럼 요약 버튼 클릭.
		a.RespondProgress()
		a.RespondButtonActionScrumSummary()
	case ButtonActionConfig:
		// 설정 버튼 클릭.
		a.RespondProgress()
		a.RespondButtonActionConfig()
	}
}

func (a *ActionResponder) RespondSelectAction(action slack.InteractiveAction) {
	switch action.ActionID {
	case ConfigActionDailyScrumEnable:
		// 데일리 스크럼 알림 설정.
		a.RespondConfigActionDailyScrumEnable(action)
	case ConfigActionWeeklyReportEnable:
		// 주간 보고 알림 설정.
		a.RespondConfigActionWeeklyReportEnable(action)
	default:
		slog.Warn("undefined select action", slog.String("action_id", action.ActionID))
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

func (a *ActionResponder) RespondButtonActionScrumList() {
	messages := listScrumMessages(a.AppToken, a.BotToken, a.Payload.Container.ChannelID, a.Payload.User.ID)
	Respond(a.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeScrumListMessage(messages),
		ReplaceOriginal: true,
	})
}

func (a *ActionResponder) RespondButtonActionScrumSummary() {
	messages := listScrumMessages(a.AppToken, a.BotToken, a.Payload.Container.ChannelID, a.Payload.User.ID)
	Respond(a.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeScrumSummaryMessage(messages),
		ReplaceOriginal: true,
	})
}

func (a *ActionResponder) RespondButtonActionConfig() {
	Respond(a.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeConfigMessage(a.Payload.Channel.ID),
		ReplaceOriginal: true,
	})
}

func (a *ActionResponder) RespondConfigActionDailyScrumEnable(action slack.InteractiveAction) {
	_, err := enableDailyScrumConfig(
		a.Payload.Container.ChannelID,
		action.SelectedOption.Value == "on",
	)
	if err != nil {
		slog.Warn("failed to enable daily scrum config")
	}

	Respond(a.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeConfigMessage(a.Payload.Container.ChannelID),
		ReplaceOriginal: true,
	})
}

func (a *ActionResponder) RespondConfigActionWeeklyReportEnable(action slack.InteractiveAction) {
	_, err := enableWeeklyReportConfig(
		a.Payload.Container.ChannelID,
		action.SelectedOption.Value == "on",
	)
	if err != nil {
		slog.Warn("failed to enable weekly report config")
	}

	Respond(a.Payload.ResponseURL, &slack.InteractiveResponsePayload{
		Blocks:          makeConfigMessage(a.Payload.Container.ChannelID),
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

func listScrumMessages(appToken, botToken, channel, user string) map[float64]string {
	messageIDs := func() []float64 {
		client, err := channelconfig.NewClient()
		if err != nil {
			slog.Error("failed to create channel config client", slog.Any("error", err))
			return nil
		}
		defer func() {
			if err := client.Close(); err != nil {
				slog.Warn("failed to close channel config client", slog.Any("error", err))
			}
		}()
		ids, err := client.GetScrumMessageHistory(context.Background(), channel)
		if err != nil {
			slog.Error("failed to get scrum message history", slog.Any("error", err))
			return nil
		}
		// 월요일 이후의 메시지 ID만 반환한다.
		min := float64(kst.LastMonday(kst.Now()).Unix())
		result := make([]float64, 0, len(ids))
		for _, id := range ids {
			if id >= min {
				result = append(result, id)
			}
		}
		return result
	}()

	if len(messageIDs) == 0 {
		return nil
	}

	client := slack.Client{
		AppToken: appToken,
		BotToken: botToken,
	}

	// 스크럼 메시지의 스레드에서 요청한 사람이 작성한 메시지를 찾는다.
	type message struct {
		ts   float64
		text string
	}
	var (
		messages = make(map[float64]string, 5)
		ch       = make(chan message, len(messageIDs))
	)

	wg := sync.WaitGroup{}
	for _, id := range messageIDs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := &slack.ListRepliesRequest{
				Channel:   channel,
				Timestamp: id,
			}
			resp, err := client.ListReplies(context.Background(), req)
			if err != nil {
				slog.Warn("failed to list replies", slog.Any("error", err))
				return
			}

			for _, reply := range resp.Messages {
				if reply.User != user {
					continue
				}
				ch <- message{
					ts:   id,
					text: reply.Text,
				}
				return
			}
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for msg := range ch {
		messages[msg.ts] = msg.text
	}

	return messages
}
