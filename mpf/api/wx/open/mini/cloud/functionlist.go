/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 10:04
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

// 获取云函数列表
type functionList struct {
    wx.BaseWxOpen
    appId  string // 应用ID
    env    string // 环境id
    offset int    // 偏移量
    limit  int    // 每页限制
}

func (fl *functionList) SetEnv(env string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, env)
    if match {
        fl.env = env
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "环境id不合法", nil))
    }
}

func (fl *functionList) SetOffset(offset int) {
    if offset >= 0 {
        fl.offset = offset
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "偏移量不合法", nil))
    }
}

func (fl *functionList) SetLimit(limit int) {
    if (limit > 0) && (limit <= 20) {
        fl.limit = limit
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "每页限制不合法", nil))
    }
}

func (fl *functionList) checkData() {
    if len(fl.env) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "环境id不能为空", nil))
    }
}

func (fl *functionList) SendRequest() api.ApiResult {
    fl.checkData()

    reqData := make(map[string]interface{})
    reqData["env"] = fl.env
    reqData["offset"] = fl.offset
    reqData["limit"] = fl.limit
    reqBody := mpf.JsonMarshal(fl.ReqData)
    fl.ReqUrl = "https://api.weixin.qq.com/tcb/listfunctions?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(fl.appId)
    client, req := fl.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := fl.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewFunctionList(appId string) *functionList {
    fl := &functionList{wx.NewBaseWxOpen(), "", "", 0, 0}
    fl.appId = appId
    fl.offset = 0
    fl.limit = 10
    fl.ReqContentType = project.HTTPContentTypeJSON
    fl.ReqMethod = fasthttp.MethodPost
    return fl
}
