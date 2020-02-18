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

// 删除菜单
type menuDelete struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
}

func (md *menuDelete) SendRequest() api.ApiResult {
    md.ReqData["access_token"] = wx.NewUtilWx().GetCorpAccessToken(md.corpId, md.agentTag)
    md.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/menu/delete?" + mpf.HTTPCreateParams(md.ReqData, "none", 1)
    client, req := md.GetRequest()
    resp, result := md.SendInner(client, req, errorcode.WxCorpRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewMenuDelete(corpId, agentTag string) *menuDelete {
    agentInfo := wx.NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    md := &menuDelete{wx.NewBaseWxCorp(), "", ""}
    md.corpId = corpId
    md.agentTag = agentTag
    md.ReqData["agentid"] = agentInfo["id"]
    return md
}
