package distributed_lock

import (
	"context"
	_ "embed"
	"github.com/go-redis/redis/v8"
	"time"
)

//go:embed res/compare_and_swap.lua
var compareAndSwapScript string

//go:embed res/compare_and_swap_ex.lua
var compareAndSwapExScript string

//go:embed res/compare_and_swap_px.lua
var compareAndSwapPxScript string

//go:embed res/compare_and_swap_keep.lua
var compareAndSwapKeepScript string

//go:embed res/compare_and_delete.lua
var compareAndDelScript string

func Cas(ctx context.Context, cli *redis.Client, key string, oldVal, newVal interface{}) (bool, error) {
	res, err := cli.Eval(ctx, compareAndSwapScript, []string{key}, oldVal, newVal).Result()
	if err != nil {
		return false, err
	}
	if res == "OK" {
		return true, nil
	}
	return false, nil
}

func CasEx(ctx context.Context, cli *redis.Client, key string, oldVal, newVal interface{}, exp time.Duration) (bool, error) {
	if exp == 0 {
		return Cas(ctx, cli, key, oldVal, newVal)
	}

	var err error
	var res interface{}
	if usePrecise(exp) {
		res, err = cli.Eval(ctx, compareAndSwapPxScript, []string{key}, oldVal, newVal, formatMs(exp)).Result()
	} else if exp > 0 {
		res, err = cli.Eval(ctx, compareAndSwapExScript, []string{key}, oldVal, newVal, formatSec(exp)).Result()
	} else {
		res, err = cli.Eval(ctx, compareAndSwapKeepScript, []string{key}, oldVal, newVal).Result()
	}

	if err != nil {
		return false, err
	}
	return res == "OK", nil
}

func Cad(ctx context.Context, cli *redis.Client, key string, val interface{}) (bool, error) {
	res, err := cli.Eval(ctx, compareAndDelScript, []string{key}, val).Result()
	if err != nil {
		return false, err
	}
	return res != 0, nil
}

// 判断给的时间是否存在小数点
func usePrecise(dur time.Duration) bool {
	return dur < time.Second || dur%time.Second != 0
}

// 小于1毫秒的按照1毫秒算
func formatMs(dur time.Duration) int64 {
	if dur > 0 && dur < time.Millisecond {
		return 1
	}
	return int64(dur / time.Millisecond)
}

// 小于1秒的按照1秒算
func formatSec(dur time.Duration) int64 {
	if dur > 0 && dur < time.Second {
		return 1
	}
	return int64(dur / time.Second)
}
