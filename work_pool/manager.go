package work_pool

import "context"

func Init() {
	initGoroutinePool()
	initUnFollowHandler()
}

func GoCtx(ctx context.Context) context.Context {
	return context.WithValue(context.TODO(), "t", ctx.Value("t"))
}
