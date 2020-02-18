/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 9:53
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

// 获取云函数下载地址
type functionDownload struct {
    wx.BaseWxOpen
    appId        string // 应用ID
    env          string // 环境id
    functionName string // 云函数名
}

func (fd *functionDownload) SetEnv(env string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, env)
    if match {
        fd.env = env
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "环境id不合法", nil))
    }
}

func (fd *functionDownload) SetFunctionName(functionName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, functionName)
    if match {
        fd.functionName = functionName
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "云函数名不合法", nil))
    }
}

func (fd *functionDownload) checkData() {
    if len(fd.env) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "环境id不能为空", nil))
    }
    if len(fd.functionName) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "云函数名不能为空", nil))
    }
    fd.ReqData["env"] = fd.env
    fd.ReqData["function_name"] = fd.functionName
}

func (fd *functionDownload) SendRequest() api.APIResult {
    fd.checkData()

    reqBody := mpf.JSONMarshal(fd.ReqData)
    fd.ReqURI = "https://api.weixin.qq.com/tcb/downloadfunction?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(fd.appId)
    client, req := fd.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := fd.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewFunctionDownload(appId string) *functionDownload {
    fd := &functionDownload{wx.NewBaseWxOpen(), "", "", ""}
    fd.appId = appId
    fd.ReqContentType = project.HTTPContentTypeJSON
    fd.ReqMethod = fasthttp.MethodPost
    return fd
}
