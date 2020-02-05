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

// 删除标签
type tagDelete struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    tagId    string // 标签id
}

func (td *tagDelete) SetTagId(tagId string) {
    if len(tagId) > 0 {
        td.tagId = tagId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "标签id不合法", nil))
    }
}

func (td *tagDelete) checkData() {
    if len(td.tagId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "标签id不能为空", nil))
    }
    td.ReqData["tagid"] = td.tagId
}

func (td *tagDelete) SendRequest(getType string) api.ApiResult {
    td.checkData()

    td.ReqData["access_token"] = wx.NewUtilWx().GetCorpCache(td.corpId, td.agentTag, getType)
    td.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/tag/delete?" + mpf.HttpCreateParams(td.ReqData, "none", 1)
    client, req := td.GetRequest()

    resp, result := td.SendInner(client, req, errorcode.WxCorpRequestGet)
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

func NewTagDelete(corpId, agentTag string) *tagDelete {
    td := &tagDelete{wx.NewBaseWxCorp(), "", "", ""}
    td.corpId = corpId
    td.agentTag = agentTag
    return td
}
