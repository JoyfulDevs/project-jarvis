package app

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/joyfuldevs/project-jarvis/pkg/kst"
	"github.com/joyfuldevs/project-jarvis/pkg/slack/blockkit"
	aigateway "github.com/joyfuldevs/project-jarvis/service/aigateway/client"
	channelconfig "github.com/joyfuldevs/project-jarvis/service/channelconfig/client"
	dataportal "github.com/joyfuldevs/project-jarvis/service/dataportal/client"
)

func makeProgressMessage() []blockkit.SlackBlock {
	return []blockkit.SlackBlock{
		&blockkit.SectionBlock{
			Text: &blockkit.TextObject{
				Type:  blockkit.TextTypePlainText,
				Text:  "⏱️ 처리 중입니다. 잠시만 기다려 주세요.",
				Emoji: true,
			},
		},
	}
}

func makeConfigMessage(channel string) []blockkit.SlackBlock {
	client, err := channelconfig.NewClient()
	if err != nil {
		return makeErrorMessage(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.Error("failed to close client", slog.Any("error", err))
		}
	}()

	config, err := client.GetChannelConfig(context.Background(), channel)
	if err != nil {
		return makeErrorMessage(err)
	}

	blocks := make([]blockkit.SlackBlock, 0, 8)
	blocks = append(blocks,
		&blockkit.HeaderBlock{
			Text: blockkit.TextObject{
				Type:  blockkit.TextTypePlainText,
				Text:  "⭐️ Jarvis 설정",
				Emoji: true,
			},
		},
		&blockkit.DividerBlock{},
		makeConfigOnOffSection(
			"스크럼 알림 설정",
			ConfigActionDailyScrumEnable,
			config.DailyScrum.Enabled,
		),
		makeConfigOnOffSection(
			"주간보고 알림 설정",
			ConfigActionWeeklyReportEnable,
			config.WeeklyReport.Enabled,
		),
		makeDoneButton(),
	)
	return blocks
}

func makeConfigOnOffSection(title string, actionID string, on bool) *blockkit.SectionBlock {
	return &blockkit.SectionBlock{
		Text: &blockkit.TextObject{
			Type: blockkit.TextTypeMarkdown,
			Text: title,
		},
		Accessory: &blockkit.SelectElement{
			ActionID: actionID,
			Type:     blockkit.ElementTypeSelect,
			Placeholder: blockkit.TextObject{
				Type:  blockkit.TextTypePlainText,
				Text:  map[bool]string{true: "ON", false: "OFF"}[on],
				Emoji: true,
			},
			Options: []blockkit.OptionBlockObject{
				{
					Text: blockkit.TextObject{
						Type:  blockkit.TextTypePlainText,
						Text:  "ON",
						Emoji: true,
					},
					Value: "on",
				},
				{
					Text: blockkit.TextObject{
						Type:  blockkit.TextTypePlainText,
						Text:  "OFF",
						Emoji: true,
					},
					Value: "off",
				},
			},
		},
	}
}

func makeErrorMessage(err error) []blockkit.SlackBlock {
	blocks := make([]blockkit.SlackBlock, 0, 4)
	blocks = append(blocks,
		&blockkit.HeaderBlock{
			Text: blockkit.TextObject{
				Type:  blockkit.TextTypePlainText,
				Text:  "⚠️ 에러! 아래 내용을 개발자에게 공유해주세요.",
				Emoji: true,
			},
		},
		&blockkit.SectionBlock{
			Text: &blockkit.TextObject{
				Type: blockkit.TextTypeMarkdown,
				Text: fmt.Sprintf("```%s```", err.Error()),
			},
		},
	)

	return blocks
}

func makeGuideMessage() []blockkit.SlackBlock {
	blocks := make([]blockkit.SlackBlock, 0, 4)
	blocks = append(blocks,
		&blockkit.HeaderBlock{
			Text: blockkit.TextObject{
				Type:  blockkit.TextTypePlainText,
				Text:  "🚫 잘못된 명령어 입니다.",
				Emoji: true,
			},
		},
		&blockkit.DividerBlock{},
		&blockkit.ActionBlock{
			Elements: []blockkit.SlackBlockElement{
				&blockkit.ButtonElement{
					ActionID: ButtonActionManual,
					Text: blockkit.TextObject{
						Type:  blockkit.TextTypePlainText,
						Text:  "📋 지원 기능 보기",
						Emoji: true,
					},
				},
				&blockkit.ButtonElement{
					ActionID: ButtonActionDone,
					Text: blockkit.TextObject{
						Type:  blockkit.TextTypePlainText,
						Text:  "✅ 완료",
						Emoji: true,
					},
				},
			},
		},
	)

	return blocks
}

func makeManualMessage() []blockkit.SlackBlock {
	return []blockkit.SlackBlock{
		&blockkit.HeaderBlock{
			Text: blockkit.TextObject{
				Type:  blockkit.TextTypePlainText,
				Text:  "⭐️ 지원 기능",
				Emoji: true,
			},
		},
		&blockkit.DividerBlock{},
		&blockkit.SectionBlock{
			Text: &blockkit.TextObject{
				Type: blockkit.TextTypeMarkdown,
				Text: fmt.Sprintf("*공휴일 안내*    `/자비스 %s`", CommandHolidayCalendar),
			},
			Accessory: &blockkit.ButtonElement{
				ActionID: ButtonActionHolidayCalendar,
				Text: blockkit.TextObject{
					Type: blockkit.TextTypePlainText,
					Text: "실행",
				},
			},
		},
		&blockkit.DividerBlock{},
		&blockkit.SectionBlock{
			Text: &blockkit.TextObject{
				Type: blockkit.TextTypeMarkdown,
				Text: fmt.Sprintf("*초단기 날씨 예보*    `/자비스 %s`", CommandForecast),
			},
			Accessory: &blockkit.ButtonElement{
				ActionID: ButtonActionForecast,
				Text: blockkit.TextObject{
					Type: blockkit.TextTypePlainText,
					Text: "실행",
				},
			},
		},
		&blockkit.DividerBlock{},
		&blockkit.SectionBlock{
			Text: &blockkit.TextObject{
				Type: blockkit.TextTypeMarkdown,
				Text: fmt.Sprintf("*스크럼 작성 이력*    `/자비스 %s`", CommandScrumList),
			},
			Accessory: &blockkit.ButtonElement{
				ActionID: ButtonActionScrumList,
				Text: blockkit.TextObject{
					Type: blockkit.TextTypePlainText,
					Text: "실행",
				},
			},
		},
		&blockkit.DividerBlock{},
		&blockkit.SectionBlock{
			Text: &blockkit.TextObject{
				Type: blockkit.TextTypeMarkdown,
				Text: fmt.Sprintf("*스크럼 요약*    `/자비스 %s`", CommandScrumSummary),
			},
			Accessory: &blockkit.ButtonElement{
				ActionID: ButtonActionScrumSummary,
				Text: blockkit.TextObject{
					Type: blockkit.TextTypePlainText,
					Text: "실행",
				},
			},
		},
		&blockkit.DividerBlock{},
		&blockkit.SectionBlock{
			Text: &blockkit.TextObject{
				Type: blockkit.TextTypeMarkdown,
				Text: fmt.Sprintf("*채널 알림 설정*    `/자비스 %s`", CommandConfig),
			},
			Accessory: &blockkit.ButtonElement{
				ActionID: ButtonActionConfig,
				Text: blockkit.TextObject{
					Type: blockkit.TextTypePlainText,
					Text: "실행",
				},
			},
		},
		&blockkit.DividerBlock{},
		makeDoneButton(),
	}
}

func makeHolidayCalendarMessage() []blockkit.SlackBlock {
	getCalendar := func(year, month int) map[int]string {
		client, err := dataportal.NewClient()
		if err != nil {
			slog.Error("failed to create client", slog.Any("error", err))
			return nil
		}
		defer func() { _ = client.Close() }()

		holidays, err := client.ListHolidays(context.Background(), year, int(month))
		if err != nil {
			slog.Error("failed to get holiday calendar", slog.Any("error", err))
			return nil
		}

		calendar := make(map[int]string, len(holidays))
		for _, holiday := range holidays {
			calendar[int(holiday.Day)] = holiday.Name
		}

		return calendar
	}

	makeField := func(t time.Time) *blockkit.TextObject {
		year, month, _ := t.Date()
		return makeHolidayField(year, int(month), getCalendar(year, int(month)))
	}
	fields := make([]blockkit.TextObject, 0, 2)
	fields = append(fields, *makeField(kst.Now()))
	fields = append(fields, *makeField(kst.Now().AddDate(0, 1, 0)))
	return []blockkit.SlackBlock{
		&blockkit.HeaderBlock{
			Text: blockkit.TextObject{
				Type:  blockkit.TextTypePlainText,
				Text:  "🗓️ 공휴일 안내",
				Emoji: true,
			},
		},
		&blockkit.DividerBlock{},
		&blockkit.SectionBlock{
			Fields: fields,
		},
		&blockkit.DividerBlock{},
		makeDoneButton(),
	}
}

func makeHolidayField(year, month int, calendar map[int]string) *blockkit.TextObject {
	days := make([]int, 0, 10)
	for day := range calendar {
		days = append(days, day)
	}
	sort.Ints(days)

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("🗓️ *%d년 %d월 공휴일 목록*\n\n", year, month))

	for _, day := range days {
		weekday := kst.Weekday(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC))
		builder.WriteString(fmt.Sprintf("%02d일 (%s요일) %s\n", day, weekday, calendar[day]))
	}

	if len(days) == 0 {
		builder.WriteString("공휴일이 없어요 😥")
	}

	return &blockkit.TextObject{
		Type: blockkit.TextTypeMarkdown,
		Text: builder.String(),
	}
}

func makeScrumListMessage(messages map[float64]string) []blockkit.SlackBlock {
	var (
		keys   = make([]float64, 0, len(messages))
		fields = make([]blockkit.TextObject, 0, 5)
		blocks = make([]blockkit.SlackBlock, 0, 8)
	)

	for key := range messages {
		keys = append(keys, key)
	}
	sort.Float64s(keys)

	for _, key := range keys {
		message := messages[key]
		builder := strings.Builder{}
		builder.WriteString("`")
		t := kst.KST(time.UnixMilli(int64(key * 1000)))
		builder.WriteString(kst.Weekday(t))
		builder.WriteString("요일")
		builder.WriteString("`\n")
		builder.WriteString("```")
		builder.WriteString(message)
		builder.WriteString("```")
		fields = append(fields, blockkit.TextObject{
			Type: blockkit.TextTypeMarkdown,
			Text: builder.String(),
		})
	}

	blocks = append(blocks,
		&blockkit.HeaderBlock{
			Text: blockkit.TextObject{
				Type:  blockkit.TextTypePlainText,
				Text:  "📝 스크럼 작성 이력",
				Emoji: true,
			},
		},
		&blockkit.DividerBlock{},
		&blockkit.SectionBlock{
			Fields: fields,
		},
		&blockkit.DividerBlock{},
		makeDoneButton(),
	)

	return blocks
}

func makeScrumSummaryMessage(messages map[float64]string) []blockkit.SlackBlock {
	texts := make([]string, 0, len(messages)+1)
	texts = append(texts, "다음 내용들을 바탕으로 이번주에 작업한 내용을 정확하게 요약해주세요.")
	for ts, msg := range messages {
		weekday := kst.Weekday(time.UnixMilli(int64(ts * 1000)))
		texts = append(texts, weekday+"요일 = "+msg)
	}
	client, err := aigateway.NewClient()
	if err != nil {
		return makeErrorMessage(err)
	}
	result, err := client.GenerateText(context.Background(), texts)
	if err != nil {
		slog.Error("failed to generate text", slog.Any("error", err))
		return makeErrorMessage(err)
	}
	return []blockkit.SlackBlock{
		&blockkit.HeaderBlock{
			Text: blockkit.TextObject{
				Type:  blockkit.TextTypePlainText,
				Text:  "🤖 스크럼 요약",
				Emoji: true,
			},
		},
		&blockkit.DividerBlock{},
		&blockkit.SectionBlock{
			Text: &blockkit.TextObject{
				Type: blockkit.TextTypeMarkdown,
				Text: fmt.Sprintf("```%s```", result),
			},
		},
		&blockkit.DividerBlock{},
		makeDoneButton(),
	}
}

func makeForecastMessage() []blockkit.SlackBlock {
	client, err := dataportal.NewClient()
	if err != nil {
		slog.Error("failed to create client", slog.Any("error", err))
		return makeErrorMessage(err)
	}

	resp, err := client.GetUltraShortTermForecast(context.Background(), 60, 123)
	if err != nil {
		slog.Error("failed to get ultra short term forecast", slog.Any("error", err))
		return makeErrorMessage(err)
	}

	blocks := make([]blockkit.SlackBlock, 0, len(resp)+4)
	blocks = append(blocks,
		&blockkit.HeaderBlock{
			Text: blockkit.TextObject{
				Type:  blockkit.TextTypePlainText,
				Text:  "🌤️ 날씨 정보",
				Emoji: true,
			},
		},
		&blockkit.DividerBlock{},
	)

	for _, item := range resp {
		t, err := time.Parse("200601021504", item.Time)
		if err != nil {
			slog.Warn("failed to parse time", slog.Any("error", err), slog.String("time", item.Time))
			continue
		}
		blocks = append(blocks, &blockkit.SectionBlock{
			Text: &blockkit.TextObject{
				Type: blockkit.TextTypeMarkdown,
				Text: fmt.Sprintf(
					"*%s* / 기온: `%d°C`, 하늘 상태: `%s`, 강수 형태: `%s`\n",
					t.Format("15:04"),
					item.Temperature,
					item.Sky,
					item.Precipitation,
				),
			},
		})
	}

	blocks = append(blocks,
		&blockkit.DividerBlock{},
		makeDoneButton(),
	)

	return blocks
}

func makeDoneButton() blockkit.SlackBlock {
	return &blockkit.ActionBlock{
		Elements: []blockkit.SlackBlockElement{
			&blockkit.ButtonElement{
				ActionID: ButtonActionDone,
				Text: blockkit.TextObject{
					Type:  blockkit.TextTypePlainText,
					Text:  "✅ 완료",
					Emoji: true,
				},
			},
		},
	}
}
