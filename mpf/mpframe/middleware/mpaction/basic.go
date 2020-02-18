/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 12:48
 */
package mpaction

import (
    "regexp"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpresp"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/kataras/iris/v12/context"
)

// 动作日志
func NewBasicLog() context.Handler {
    return func(ctx context.Context) {
        reqUrl := ctx.Values().GetString(project.DataParamKeyReqURL)
        mplog.LogInfo(reqUrl + " action-enter")

        // 业务结束日志
        actionStart := time.Now()
        defer func() {
            costTime := time.Since(actionStart).Seconds()
            costTimeStr := strconv.FormatFloat(costTime, 'f', 6, 64)
            mplog.LogInfo(reqUrl + " action-exit,cost_time: " + costTimeStr + "s")
            if costTime >= ctx.Application().ConfigurationReadOnly().GetOther()["timeout_action"].(float64) {
                mplog.LogWarn("handle " + reqUrl + " action-timeout,cost_time: " + costTimeStr + "s")
            }
        }()
    }
}

// 简单签名验证(只有api模块验证签名)
func NewBasicSignSimple() context.Handler {
    return func(ctx context.Context) {
        errMsg := ""
        sign := ctx.URLParamDefault("_sign", "")
        if len(sign) == 0 {
            errMsg = "签名不能为空"
        } else if match, _ := regexp.MatchString(`^[0-9]{10}[0-9a-z]+$`, sign); !match {
            errMsg = "签名不合法"
        } else {
            nowTime := time.Now().Unix()
            signTime, _ := strconv.ParseInt(sign[0:10], 10, 64)
            leftTime := nowTime - signTime
            if (leftTime > 30) && (leftTime < -30) {
                errMsg = "签名不正确"
            } else {
                signNew := ""
                conf := mpf.NewConfig().GetConfig("mpauth")
                secret := conf.GetString("simple." + mpf.EnvType() + mpf.EnvProjectTag() + ".secret")
                hashMethod := conf.GetString("simple." + mpf.EnvType() + mpf.EnvProjectTag() + ".hash")
                switch hashMethod {
                case "md5":
                    signNew = mpf.HashMd5(sign[0:10]+secret, "")
                case "sha1":
                    signNew = mpf.HashSha1(sign[0:10], secret)
                case "sha256":
                    signNew = mpf.HashSha256(sign[0:10], secret)
                case "sha512":
                    signNew = mpf.HashSha512(sign[0:10], secret)
                default:
                    signNew = mpf.HashMd5(sign[0:10]+secret, "")
                }
                if signNew != sign[10:] {
                    errMsg = "签名不正确"
                }
            }
        }

        if len(errMsg) > 0 {
            result := mpresponse.NewResultProblem()
            result.Tag = "validator-sign"
            result.Title = "签名错误"
            result.Code = errorcode.CommonValidatorSign
            result.Msg = errMsg
            ctx.Problem(mpresp.GetProblemHandleBasic(result, 30*time.Second))
            mpresp.NewBasicEnd()(ctx)
        } else {
            ctx.Next()
        }
    }
}
