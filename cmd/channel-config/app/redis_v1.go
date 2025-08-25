package app

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/devafterdark/project-jarvis/service/channelconfig/server"
)

var _ server.ServiceV1 = (*RedisService)(nil)

func (r *RedisService) Subscribe(ctx context.Context, channel string, feature server.FeatureV1) ([]server.FeatureV1, error) {
	key := fmt.Sprintf("subscription:%s", channel)

	flag, err := r.rdb.Get(ctx, key).Uint64()
	if err == redis.Nil {
		flag = 0
	} else if err != nil {
		return nil, err
	}

	new := flag | (1 << uint64(feature))
	err = r.rdb.Set(ctx, key, new, 0).Err()
	if err != nil {
		return nil, err
	}

	r.rdb.SAdd(ctx, fmt.Sprintf("channels:%d", feature), channel)

	return flagToFeatures(new), nil
}

func (r *RedisService) Unsubscribe(ctx context.Context, channel string, feature server.FeatureV1) ([]server.FeatureV1, error) {
	key := fmt.Sprintf("subscription:%s", channel)

	flag, err := r.rdb.Get(ctx, key).Uint64()
	if err == redis.Nil {
		flag = 0
	} else if err != nil {
		return nil, err
	}

	new := flag &^ (1 << uint64(feature))
	err = r.rdb.Set(ctx, key, new, 0).Err()
	if err != nil {
		return nil, err
	}

	r.rdb.SRem(ctx, fmt.Sprintf("channels:%d", feature), channel)

	return flagToFeatures(new), nil
}

func (r *RedisService) UnsubscribeAll(ctx context.Context, channel string) error {
	key := fmt.Sprintf("subscription:%s", channel)
	if err := r.rdb.Set(ctx, key, 0, 0).Err(); err != nil {
		return err
	}

	keys, err := r.rdb.Keys(ctx, "channels:*").Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}

	for _, key := range keys {
		r.rdb.SRem(ctx, key, channel)
	}

	return nil
}

func (r *RedisService) ListSubscriptions(ctx context.Context, channel string) ([]server.FeatureV1, error) {
	key := fmt.Sprintf("subscription:%s", channel)

	flag, err := r.rdb.Get(ctx, key).Uint64()
	if err == redis.Nil {
		flag = 0
	} else if err != nil {
		return nil, err
	}

	return flagToFeatures(flag), nil
}

func (r *RedisService) ListChannelsByFeature(ctx context.Context, feature server.FeatureV1) ([]string, error) {
	key := fmt.Sprintf("channels:%d", feature)
	channels, err := r.rdb.SMembers(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return channels, nil
}

func flagToFeatures(flag uint64) []server.FeatureV1 {
	var features []server.FeatureV1
	for i := range 16 {
		if flag&(1<<i) > 0 {
			features = append(features, server.FeatureV1(i))
		}
	}
	return features
}
