package mpspecial

import "github.com/kataras/iris/v12/context"

// NewBasicEmpty 空操作处理,可用于OPTIONS请求响应等特殊场景
func NewBasicEmpty() context.Handler {
    return func(ctx context.Context) {
        ctx.Next()
    }
}
