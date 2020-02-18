/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/8 0008
 * Time: 16:04
 */
package tag

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建标签
type tagCreate struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    tagName  string // 名称
    tagId    string // 标签ID
}

func (tc *tagCreate) SetTagName(tagName string) {
    if len(tagName) > 0 {
        trueName := []rune(tagName)
        tc.tagName = string(trueName[:16])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "名称不合法", nil))
    }
}

func (tc *tagCreate) SetTagId(tagId string) {
    if len(tagId) > 0 {
        tc.ReqData["tagid"] = tagId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "标签id不合法", nil))
    }
}

func (tc *tagCreate) checkData() {
    if len(tc.tagName) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "名称不能为空", nil))
    }
    tc.ReqData["tagname"] = tc.tagName
}

func (tc *tagCreate) SendRequest(getType string) api.ApiResult {
    tc.checkData()

    reqBody := mpf.JSONMarshal(tc.ReqData)
    tc.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/tag/create?access_token=" + wx.NewUtilWx().GetCorpCache(tc.corpId, tc.agentTag, getType)
    client, req := tc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := tc.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewTagCreate(corpId, agentTag string) *tagCreate {
    tc := &tagCreate{wx.NewBaseWxCorp(), "", "", "", ""}
    tc.corpId = corpId
    tc.agentTag = agentTag
    tc.ReqContentType = project.HTTPContentTypeJSON
    tc.ReqMethod = fasthttp.MethodPost
    return tc
}
