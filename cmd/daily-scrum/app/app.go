package app

import (
	"log/slog"

	"github.com/genians/endpoint-lab-slack-bot/pkg/kst"
)

func Run() {
	koreaTime := kst.Now()
	if IsHoliday(koreaTime) {
		return
	}

	// 봇이 초대된 채널 목록을 가져온다.
	channels := ListInvitedChannels()
	if len(channels) == 0 {
		slog.Info("no channel to send scrum message")
		return
	}

	// 데일리 스크럼 기능을 활성화한 채널만 필터링한다.
	targets := make([]string, 0, len(channels))
	for _, channel := range channels {
		config, err := DailyScrumConfig(channel)
		if err != nil {
			continue
		}
		if config.Enabled {
			targets = append(targets, channel)
		}
	}

	// 기능이 활성화된 채널에 스크럼 메시지를 보낸다.
	SendScrumMessage(targets, DailyScrumMessage(koreaTime))
}
