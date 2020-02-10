/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 12:47
 */
package mpresp

import (
    "os"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
)

// 获取问题处理调用链
func NewBasicHandlersProblem(problem *mpresponse.ResultProblem, retryAfter interface{}) []context.Handler {
    handlers := make([]context.Handler, 2)

    irisProblem := iris.NewProblem().
        Type("/error/"+problem.Type).
        Title(problem.Title).
        Detail(problem.Detail).
        Status(problem.Status).
        Key("req_id", problem.ReqId).
        Key("code", problem.Code).
        Key("time", problem.Time).
        Key("msg", problem.Msg)
    handlers = append(handlers, func(ctx context.Context) {
        ctx.Problem(irisProblem, iris.ProblemOptions{
            JSON:       context.JSON{Indent: ""},
            RetryAfter: retryAfter,
        })
        ctx.Next()
    })
    handlers = append(handlers, NewBasicEnd())
    return handlers
}

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
                result := mpresponse.NewResultApi()
                result.Data = data.(interface{})
                ctx.Header(project.HttpHeadKeyContentType, project.HttpContentTypeJson)
                ctx.WriteString(mpf.JsonMarshal(result))
            }

            ctx.Next()
        } else {
            problem := mpresponse.NewResultProblem()
            problem.Type = "response-empty"
            problem.Title = "响应错误"
            problem.Detail = "响应数据未设置"
            problem.Code = errorcode.CommonResponseEmpty
            problem.Msg = "响应数据不能为空"
            ctx.Do(NewBasicHandlersProblem(problem, 30*time.Second))
        }
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
