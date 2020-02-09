/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 12:47
 */
package mpresp

import (
    "os"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/kataras/iris/v12/context"
)

// 发送响应数据
func NewBasicSend() context.Handler {
    return func(ctx context.Context) {
        respData, ok := ctx.Values().GetEntry(project.DataParamKeyRespData)
        if ok {
            data := respData.Value()
            switch data.(type) {
            case string:
                ctx.Header(project.HttpHeadKeyContentType, project.HttpContentTypeText)
                ctx.WriteString(data.(string))
            default:
                result := mpresponse.NewResultBasic()
                result.Data = data.(interface{})
                ctx.Header(project.HttpHeadKeyContentType, project.HttpContentTypeJson)
                ctx.WriteString(mpf.JsonMarshal(result))
            }
        } else {
            result := mpresponse.NewResultBasic()
            result.Code = errorcode.CommonBaseServer
            result.Msg = "响应数据不能为空"
            ctx.Header(project.HttpHeadKeyContentType, project.HttpContentTypeJson)
            ctx.WriteString(mpf.JsonMarshal(result))
        }

        ctx.Next()
    }
}

// 请求响应结束
func NewBasicEnd() context.Handler {
    return func(ctx context.Context) {
        os.Unsetenv(project.DataParamKeyReqId)
        ctx.Values().Remove(project.DataParamKeyReqUrl)
        ctx.Values().Remove(project.DataParamKeyRespData)
        // 最后退出上下文的时候,不要用ctx.EndRequest(),它会导致响应的数据被复制一份
        ctx.StopExecution()
    }
}
