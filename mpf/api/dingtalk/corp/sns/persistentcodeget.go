/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 0:30
 */
package sns

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取用户的持久授权码
type persistentCodeGet struct {
    dingtalk.BaseCorp
    corpId      string
    tmpAuthCode string // 临时授权码
}

func (pcg *persistentCodeGet) SetTmpAuthCode(tmpAuthCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, tmpAuthCode)
    if match {
        pcg.tmpAuthCode = tmpAuthCode
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "临时授权码不合法", nil))
    }
}

func (pcg *persistentCodeGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pcg.tmpAuthCode) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "临时授权码不能为空", nil))
    }
    pcg.ExtendData["tmp_auth_code"] = pcg.tmpAuthCode

    pcg.ReqUrl = dingtalk.UrlService + "/sns/get_persistent_code?access_token="
    if len(pcg.corpId) > 0 {
        pcg.ReqUrl += dingtalk.NewUtil().GetCorpSnsToken(pcg.corpId)
    } else {
        pcg.ReqUrl += dingtalk.NewUtil().GetProviderSnsToken()
    }

    reqBody := mpf.JsonMarshal(pcg.ExtendData)
    client, req := pcg.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewPersistentCodeGet(corpId string) *persistentCodeGet {
    pcg := &persistentCodeGet{dingtalk.NewCorp(), "", ""}
    pcg.corpId = corpId
    pcg.ReqContentType = project.HttpContentTypeJson
    pcg.ReqMethod = fasthttp.MethodPost
    return pcg
}
