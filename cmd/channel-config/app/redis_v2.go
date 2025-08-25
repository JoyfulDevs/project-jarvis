package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/devafterdark/project-jarvis/service/channelconfig/server"
)

const (
	// 최대 2주간의 메시지 저장.
	scrumHistorySize = 14
)

var _ server.ServiceV2 = (*RedisService)(nil)

func (r *RedisService) GetChannelConfig(
	ctx context.Context,
	channel string,
) (*server.ChannelConfigV2, error) {
	config := &server.ChannelConfigV2{}

	if c, err := r.GetDailyScrumConfig(ctx, channel); err != nil {
		return nil, err
	} else {
		config.DailyScrum = c
	}

	if c, err := r.GetWeeklyReportConfig(ctx, channel); err != nil {
		return nil, err
	} else {
		config.WeeklyReport = c
	}

	return config, nil
}

func (r *RedisService) GetDailyScrumConfig(
	ctx context.Context,
	channel string,
) (*server.DailyScrumConfigV2, error) {
	key := fmt.Sprintf("config:%s:daily_scrum", channel)
	config := &server.DailyScrumConfigV2{
		Enabled: false,
	}
	data, err := r.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return config, nil
	} else if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (r *RedisService) SetDailyScrumConfig(
	ctx context.Context,
	channel string,
	config *server.DailyScrumConfigV2,
) error {
	key := fmt.Sprintf("config:%s:daily_scrum", channel)
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	if err := r.rdb.Set(ctx, key, data, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisService) GetWeeklyReportConfig(
	ctx context.Context,
	channel string,
) (*server.WeeklyReportConfigV2, error) {
	key := fmt.Sprintf("config:%s:weekly_report", channel)
	config := &server.WeeklyReportConfigV2{
		Enabled: false,
	}
	data, err := r.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return config, nil
	} else if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (r *RedisService) SetWeeklyReportConfig(
	ctx context.Context,
	channel string,
	config *server.WeeklyReportConfigV2,
) error {
	key := fmt.Sprintf("config:%s:weekly_report", channel)
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	if err := r.rdb.Set(ctx, key, data, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisService) AddScrumMessageHistory(
	ctx context.Context,
	channelID string,
	messageID float64,
) error {
	key := fmt.Sprintf("recently_sent:daily_scrum:%s", channelID)
	if err := r.rdb.LPush(ctx, key, messageID).Err(); err != nil {
		return err
	}
	if err := r.rdb.LTrim(ctx, key, 0, scrumHistorySize).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisService) GetScrumMessageHistory(
	ctx context.Context,
	channelID string,
) ([]float64, error) {
	key := fmt.Sprintf("recently_sent:daily_scrum:%s", channelID)
	ids, err := r.rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	var result []float64
	for _, id := range ids {
		value, err := strconv.ParseFloat(id, 64)
		if err != nil {
			slog.Warn("failed to parse message ID", slog.String("id", id), slog.Any("error", err))
			continue
		}
		result = append(result, value)
	}
	return result, nil
}
