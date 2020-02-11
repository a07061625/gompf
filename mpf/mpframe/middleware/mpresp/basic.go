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

// 获取错误处理
func GetProblemHandleBasic(result *mpresponse.ResultProblem, retryAfter interface{}) (context.Problem, context.ProblemOptions) {
    problem := iris.NewProblem()
    problem.Type("/error/" + result.Type)
    problem.Title(result.Title)
    problem.Status(result.Status)
    problem.Key("req_id", result.ReqId)
    problem.Key("code", result.Code)
    problem.Key("time", result.Time)
    problem.Key("msg", result.Msg)

    return problem, iris.ProblemOptions{
        JSON:       context.JSON{Indent: "  "},
        RetryAfter: retryAfter,
    }
}

// 发送响应数据
func NewBasicSend() context.Handler {
    return func(ctx context.Context) {
        respData, ok := ctx.Values().GetEntry(project.DataParamKeyRespData)
        if ok {
            data := respData.Value()
            switch data.(type) {
            case string:
                ctx.Values().Set(project.DataParamKeyRespType, project.HttpContentTypeText)
                ctx.WriteString(data.(string))
            default:
                result := mpresponse.NewResultApi()
                result.Data = data.(interface{})
                ctx.WriteString(mpf.JsonMarshal(result))
            }

            ctx.Next()
        } else {
            result := mpresponse.NewResultProblem()
            result.Type = "response-empty"
            result.Title = "响应错误"
            result.Code = errorcode.CommonResponseEmpty
            result.Msg = "响应数据不能为空"
            ctx.Problem(GetProblemHandleBasic(result, 30*time.Second))
            NewBasicEnd()(ctx)
        }
    }
}

// 请求最终清理
func HandleEndBasic(ctx context.Context) {
    ctx.StatusCode(iris.StatusOK)
    ctx.Header("Connection", "close") // 解决大量ESTABLISHED状态请求问题
    // 设置响应数据类型
    ctx.Recorder().Header().Del(project.HttpHeadKeyContentType)
    respType, ok := ctx.Values().GetEntry(project.DataParamKeyRespType)
    if ok && (project.HttpContentTypeText == respType.Value().(string)) {
        ctx.Header(project.HttpHeadKeyContentType, project.HttpContentTypeText)
    } else {
        ctx.Header(project.HttpHeadKeyContentType, project.HttpContentTypeJson)
    }

    os.Unsetenv(project.DataParamKeyReqId)
    ctx.Values().Remove(project.DataParamKeyReqUrl)
    ctx.Values().Remove(project.DataParamKeyRespData)
    ctx.Values().Remove(project.DataParamKeyRespType)
    // 最后退出上下文的时候,不要用ctx.EndRequest(),它会导致响应的数据被复制一份
    ctx.StopExecution()
    ctx.Recorder().EndResponse()
}

// 请求响应结束
func NewBasicEnd() context.Handler {
    return func(ctx context.Context) {
        HandleEndBasic(ctx)
    }
}
