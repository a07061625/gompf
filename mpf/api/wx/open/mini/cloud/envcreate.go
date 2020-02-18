/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 9:03
 */
package cloud

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建云环境
type envCreate struct {
    wx.BaseWxOpen
    appId string // 应用ID
    env   string // 环境id
    alias string // 环境别名
}

func (ec *envCreate) SetEnv(env string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, env)
    if match {
        ec.env = env
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "环境id不合法", nil))
    }
}

func (ec *envCreate) SetAlias(alias string) {
    if len(alias) > 0 {
        ec.alias = alias
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "环境别名不合法", nil))
    }
}

func (ec *envCreate) checkData() {
    if len(ec.env) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "环境id不能为空", nil))
    }
    if len(ec.alias) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "环境别名不能为空", nil))
    }
    ec.ReqData["env"] = ec.env
    ec.ReqData["alias"] = ec.alias
}

func (ec *envCreate) SendRequest() api.ApiResult {
    ec.checkData()

    reqBody := mpf.JSONMarshal(ec.ReqData)
    ec.ReqUrl = "https://api.weixin.qq.com/tcb/createenvandresource?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(ec.appId)
    client, req := ec.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ec.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewEnvCreate(appId string) *envCreate {
    ec := &envCreate{wx.NewBaseWxOpen(), "", "", ""}
    ec.appId = appId
    ec.ReqContentType = project.HTTPContentTypeJSON
    ec.ReqMethod = fasthttp.MethodPost
    return ec
}
