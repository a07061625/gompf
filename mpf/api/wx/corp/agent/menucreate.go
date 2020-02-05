/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 17:08
 */
package agent

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建菜单
type menuCreate struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    buttons  []interface{} // 菜单列表
}

func (mc *menuCreate) SetButtons(buttons []interface{}) {
    if len(buttons) > 0 {
        mc.buttons = buttons
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "菜单列表不合法", nil))
    }
}

func (mc *menuCreate) checkData() {
    if len(mc.buttons) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "菜单列表不能为空", nil))
    }
}

func (mc *menuCreate) SendRequest() api.ApiResult {
    mc.checkData()
    reqData := make(map[string]interface{})
    reqData["button"] = mc.buttons
    reqBody := mpf.JsonMarshal(reqData)

    mc.ReqData["access_token"] = wx.NewUtilWx().GetCorpAccessToken(mc.corpId, mc.agentTag)
    mc.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/menu/create?" + mpf.HttpCreateParams(mc.ReqData, "none", 1)
    client, req := mc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := mc.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewMenuCreate(corpId, agentTag string) *menuCreate {
    agentInfo := wx.NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    mc := &menuCreate{wx.NewBaseWxCorp(), "", "", make([]interface{}, 0)}
    mc.corpId = corpId
    mc.agentTag = agentTag
    mc.ReqData["agentid"] = agentInfo["id"]
    mc.ReqContentType = project.HttpContentTypeJson
    mc.ReqMethod = fasthttp.MethodPost
    return mc
}
