// Package mpresp basic
// User: 姜伟
// Time: 2020-02-25 10:55:20
package mpresp

import (
    "os"
    "time"

    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
)

// GetProblemHandleBasic 获取错误处理
func GetProblemHandleBasic(result *mpresponse.ResultProblem, retryAfter interface{}) (context.Problem, context.ProblemOptions) {
    problem := iris.NewProblem()
    problem.Title(result.Title)
    problem.Status(result.Status)
    problem.Key("tag", result.Tag)
    problem.Key("req_id", result.ReqID)
    problem.Key("code", result.Code)
    problem.Key("time", result.Time)
    problem.Key("msg", result.Msg)

    return problem, iris.ProblemOptions{
        JSON:       context.JSON{Indent: "  "},
        RetryAfter: retryAfter,
    }
}

// NewBasicSend 发送响应数据
func NewBasicSend() context.Handler {
    return func(ctx context.Context) {
        respData, ok := ctx.Values().GetEntry(project.DataParamKeyRespData)
        if ok {
            data := respData.Value()
            switch data.(type) {
            case string:
                ctx.Values().Set(project.DataParamKeyRespType, project.HTTPContentTypeText)
                ctx.WriteString(data.(string))
            default:
                result := mpresponse.NewResultAPI()
                result.Data = data.(interface{})
                ctx.JSON(result, context.JSON{Indent: "  "})
            }

            ctx.Next()
        } else {
            result := mpresponse.NewResultProblem()
            result.Tag = "response-empty"
            result.Title = "响应错误"
            result.Code = errorcode.CommonResponseEmpty
            result.Msg = "响应数据不能为空"
            ctx.Problem(GetProblemHandleBasic(result, 30*time.Second))
            NewBasicEnd()(ctx)
        }
    }
}

// HandleEndBasic 请求最终清理
func HandleEndBasic(ctx context.Context) {
    ctx.StatusCode(iris.StatusOK)
    ctx.Header("Connection", "close") // 解决大量ESTABLISHED状态请求问题
    // 设置响应数据类型
    ctx.Recorder().Header().Del(project.HTTPHeadKeyContentType)
    respType, ok := ctx.Values().GetEntry(project.DataParamKeyRespType)
    if ok && (project.HTTPContentTypeText == respType.Value().(string)) {
        ctx.Header(project.HTTPHeadKeyContentType, project.HTTPContentTypeText)
    } else {
        ctx.Header(project.HTTPHeadKeyContentType, project.HTTPContentTypeJSON)
    }

    os.Unsetenv(project.DataParamKeyReqID)
    ctx.Values().Remove(project.DataParamKeyReqURL)
    ctx.Values().Remove(project.DataParamKeyRespData)
    ctx.Values().Remove(project.DataParamKeyRespType)
    // 最后退出上下文的时候,不要用ctx.EndRequest(),它会导致响应的数据被复制一份
    ctx.StopExecution()
    ctx.Recorder().EndResponse()
}

// NewBasicEnd 请求响应结束
func NewBasicEnd() context.Handler {
    return func(ctx context.Context) {
        HandleEndBasic(ctx)
    }
}
