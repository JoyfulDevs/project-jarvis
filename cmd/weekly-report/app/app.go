package app

import (
	"log/slog"

	"github.com/devafterdark/project-jarvis/pkg/kst"
)

func Run() {
	koreaTime := kst.Now()
	if !IsLastWorkday(koreaTime) {
		return
	}

	// 봇이 초대된 채널 목록을 가져온다.
	channels := ListInvitedChannels()
	if len(channels) == 0 {
		slog.Info("no channel to send scrum message")
		return
	}

	// 주간 업무 보고 기능을 활성화한 채널만 필터링한다.
	targets := make([]string, 0, len(channels))
	for _, channel := range channels {
		config, err := WeeklyReportConfig(channel)
		if err != nil {
			continue
		}
		if config.Enabled {
			targets = append(targets, channel)
		}
	}

	// 채널에 주간 업무 보고 메시지를 보낸다.
	SendReportMessage(targets, WeeklyReportMessage())
}
