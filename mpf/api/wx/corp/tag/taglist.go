/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/8 0008
 * Time: 23:28
 */
package tag

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 获取标签列表
type tagList struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
}

func (tl *tagList) SendRequest(getType string) api.APIResult {
    tl.ReqData["access_token"] = wx.NewUtilWx().GetCorpCache(tl.corpId, tl.agentTag, getType)
    tl.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/tag/list?" + mpf.HTTPCreateParams(tl.ReqData, "none", 1)
    client, req := tl.GetRequest()

    resp, result := tl.SendInner(client, req, errorcode.WxCorpRequestGet)
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

func NewTagList(corpId, agentTag string) *tagList {
    tl := &tagList{wx.NewBaseWxCorp(), "", ""}
    tl.corpId = corpId
    tl.agentTag = agentTag
    return tl
}
