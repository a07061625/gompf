/**
 * 版本中间件,需要设置请求头,下面两个方式任选其一
 *   Accept: "application/json; version=1.0"
 *   Accept-Version: "1.0"
 * User: 姜伟
 * Date: 2020/2/9 0009
 * Time: 23:38
 */
package mpversion

import (
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/kataras/iris/v12/context"
    "github.com/kataras/iris/v12/versioning"
)

// 版本错误
func NewBasicError() context.Handler {
    return func(ctx context.Context) {
        errMsg := ""
        apiVersion := versioning.GetVersion(ctx)
        minVersion := ctx.Application().ConfigurationReadOnly().GetOther()["version_min"].(string)
        maxVersion := ctx.Application().ConfigurationReadOnly().GetOther()["version_max"].(string)
        if apiVersion == versioning.NotFound {
            errMsg = "API版本必须填写"
        } else if versioning.Match(ctx, "< "+minVersion) {
            errMsg = "API版本已废弃"
        } else if versioning.Match(ctx, "> "+maxVersion) {
            errMsg = "API版本不支持"
        }
        if len(errMsg) > 0 {
            result := mpresponse.NewResultBasic()
            result.Code = errorcode.CommonRequestFormat
            result.Msg = errMsg
            ctx.Recorder().Header().Set(project.HttpHeadKeyContentType, project.HttpContentTypeJson)
            ctx.Recorder().SetBodyString(mpf.JsonMarshal(result))
            ctx.StopExecution()
        }
    }
}

// 版本将移除
func NewBasicDeprecated(handler context.Handler, warn string, info string) context.Handler {
    return versioning.Deprecated(handler, versioning.DeprecationOptions{
        WarnMessage:     warn,
        DeprecationDate: time.Now(),
        DeprecationInfo: info,
    })
}

// 版本匹配
func NewBasicMatcher(handlers map[string]context.Handler) context.Handler {
    return versioning.NewMatcher(versioning.Map(handlers))
}