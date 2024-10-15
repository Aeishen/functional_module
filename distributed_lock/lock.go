package distributed_lock

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

const defaultExp = 10 * time.Second // 锁默认时间

type RedisLock struct {
	cli    *redis.Client
	key    string
	cancel context.CancelFunc
}

func NewRedisLock(cli *redis.Client, key string) *RedisLock {
	return &RedisLock{
		cli: cli,
		key: key,
	}
}

func (r *RedisLock) Lock(ctx context.Context) (bool, error) {
	ok, err := r.cli.SetNX(ctx, r.key, 1, defaultExp).Result()
	if err != nil {
		return false, err
	}
	c, cancel := context.WithCancel(ctx)
	r.cancel = cancel
	go r.refresh(c)
	return ok, nil
}

func (r *RedisLock) UnLock(ctx context.Context) (bool, error) {
	ok, err := Cad(ctx, r.cli, r.key, 1)
	if err != nil {
		return false, err
	}

	if ok {
		r.cancel()
	}
	return ok, nil
}

func (r *RedisLock) refresh(ctx context.Context) {
	// 周期性的定时器一直续租
	ticker := time.NewTicker(defaultExp / 4)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.cli.Expire(ctx, r.key, defaultExp)
		}
	}
}
