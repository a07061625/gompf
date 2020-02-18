/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 0:30
 */
package sns

import (
    "regexp"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 通过临时授权码获取授权用户的个人信息
type userInfoByCode struct {
    dingtalk.BaseCorp
    corpId      string
    atType      string
    tmpAuthCode string // 临时授权码
}

func (uic *userInfoByCode) SetTmpAuthCode(tmpAuthCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, tmpAuthCode)
    if match {
        uic.tmpAuthCode = tmpAuthCode
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "临时授权码不合法", nil))
    }
}

func (uic *userInfoByCode) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(uic.tmpAuthCode) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "临时授权码不能为空", nil))
    }
    uic.ExtendData["tmp_auth_code"] = uic.tmpAuthCode

    timestamp := strconv.FormatInt(time.Now().Unix(), 10)
    if uic.atType == dingtalk.AccessTokenTypeCorp {
        conf := dingtalk.NewConfig().GetCorp(uic.corpId)
        uic.ReqData["accessKey"] = conf.GetLoginAppId()
        uic.ReqData["signature"] = dingtalk.NewUtil().CreateApiSign(timestamp, conf.GetLoginAppSecret())
    } else {
        conf := dingtalk.NewConfig().GetProvider()
        uic.ReqData["accessKey"] = conf.GetLoginAppId()
        uic.ReqData["signature"] = dingtalk.NewUtil().CreateApiSign(timestamp, conf.GetLoginAppSecret())
    }
    uic.ReqData["timestamp"] = timestamp
    uic.ReqUrl = dingtalk.UrlService + "/sns/getuserinfo_bycode?" + mpf.HttpCreateParams(uic.ReqData, "none", 1)

    reqBody := mpf.JsonMarshal(uic.ExtendData)
    client, req := uic.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewUserInfoByCode(corpId, atType string) *userInfoByCode {
    uic := &userInfoByCode{dingtalk.NewCorp(), "", "", ""}
    uic.corpId = corpId
    uic.atType = atType
    uic.ReqContentType = project.HTTPContentTypeJSON
    uic.ReqMethod = fasthttp.MethodPost
    return uic
}
