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
)

// 获取菜单
type menuGet struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
}

func (mg *menuGet) SendRequest() api.ApiResult {
    mg.ReqData["access_token"] = wx.NewUtilWx().GetCorpAccessToken(mg.corpId, mg.agentTag)
    mg.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/menu/get?" + mpf.HttpCreateParams(mg.ReqData, "none", 1)
    client, req := mg.GetRequest()
    resp, result := mg.SendInner(client, req, errorcode.WxCorpRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewMenuGet(corpId, agentTag string) *menuGet {
    agentInfo := wx.NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    mg := &menuGet{wx.NewBaseWxCorp(), "", ""}
    mg.corpId = corpId
    mg.agentTag = agentTag
    mg.ReqData["agentid"] = agentInfo["id"]
    return mg
}
