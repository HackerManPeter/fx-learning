package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hackermanpeter/fx-learning/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Cache struct {
	cfg    *config.Config
	logger *zap.Logger
	rdb    *redis.Client
}

func NewCache(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) *Cache {
	addr := fmt.Sprintf("%v:%v", cfg.Cache.Host, cfg.Cache.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Sugar().Infow("closing redis")
			return rdb.Close()
		},
	})

	return &Cache{
		cfg,
		logger,
		rdb,
	}
}

func (c *Cache) Get(ctx context.Context, key string) (any, error) {
	c.logger.Sugar().Infow("🚀 ~ Cache ~ c.Get ~ Getting value from redis", "key", key)
	value, err := c.rdb.Get(ctx, c.generateKey(key)).Result()
	if err != nil && err != redis.Nil {
		c.logger.Sugar().Errorw("🚀 ~ Cache ~ c.Get ~ unable to read from redis", "err", err, "key", key)
		return nil, err
	} else if err != nil && err == redis.Nil {
		return nil, nil
	}

	return value, nil
}

func (c *Cache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	// var v string
	// v, ok := value.(string)
	// if !ok {
	// 	vBytes, err := json.Marshal(v)
	// 	if err != nil {
	// 		c.logger.Sugar().Errorw("🚀 ~ Cache ~ c.Set ~ unable to marshal JSON", "err", err, "key", key)
	// 		return errors.New("unable to store value in redis")
	// 	}

	// 	v = string(vBytes)
	// }

	if expiration == 0 {
		expiration = 10 * time.Minute
	}

	_, err := c.rdb.Set(ctx, c.generateKey(key), value, expiration).Result()
	if err != nil {
		c.logger.Sugar().Errorw("🚀 ~ Cache ~ c.Set ~ unable to store value:", "err", err, "key", key)
		return errors.New("unable to store value")
	}

	return nil

}

func (c *Cache) Del(ctx context.Context, key string) error {
	_, err := c.rdb.Del(ctx, c.generateKey(key)).Result()
	if err != nil {
		c.logger.Sugar().Errorw("🚀 ~ Cache ~ c.Set ~ unable to delete value", "err", err, "key", key)
		return errors.New("unable to delete value")
	}

	return nil
}

func (c *Cache) generateKey(key string) string {
	return fmt.Sprintf("%v:%v", c.cfg.Cache.Prefix, key)
}
