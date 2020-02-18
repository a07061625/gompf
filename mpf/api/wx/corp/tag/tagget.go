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
    "github.com/a07061625/gompf/mpf/mperr"
)

// 获取标签成员
type tagGet struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    tagId    string // 标签id
}

func (tg *tagGet) SetTagId(tagId string) {
    if len(tagId) > 0 {
        tg.tagId = tagId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "标签id不合法", nil))
    }
}

func (tg *tagGet) checkData() {
    if len(tg.tagId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "标签id不能为空", nil))
    }
    tg.ReqData["tagid"] = tg.tagId
}

func (tg *tagGet) SendRequest(getType string) api.ApiResult {
    tg.checkData()

    tg.ReqData["access_token"] = wx.NewUtilWx().GetCorpCache(tg.corpId, tg.agentTag, getType)
    tg.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/tag/get?" + mpf.HTTPCreateParams(tg.ReqData, "none", 1)
    client, req := tg.GetRequest()

    resp, result := tg.SendInner(client, req, errorcode.WxCorpRequestGet)
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

func NewTagGet(corpId, agentTag string) *tagGet {
    tg := &tagGet{wx.NewBaseWxCorp(), "", "", ""}
    tg.corpId = corpId
    tg.agentTag = agentTag
    return tg
}
