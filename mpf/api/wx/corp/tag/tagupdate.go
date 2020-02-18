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

// 更新标签名字
type tagUpdate struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    tagName  string // 名称
    tagId    string // 标签ID
}

func (tu *tagUpdate) SetTagName(tagName string) {
    if len(tagName) > 0 {
        trueName := []rune(tagName)
        tu.tagName = string(trueName[:16])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "名称不合法", nil))
    }
}

func (tu *tagUpdate) SetTagId(tagId string) {
    if len(tagId) > 0 {
        tu.tagId = tagId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "标签id不合法", nil))
    }
}

func (tu *tagUpdate) checkData() {
    if len(tu.tagId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "标签id不能为空", nil))
    }
    if len(tu.tagName) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "名称不能为空", nil))
    }
    tu.ReqData["tagid"] = tu.tagId
    tu.ReqData["tagname"] = tu.tagName
}

func (tu *tagUpdate) SendRequest(getType string) api.ApiResult {
    tu.checkData()

    reqBody := mpf.JsonMarshal(tu.ReqData)
    tu.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/tag/update?access_token=" + wx.NewUtilWx().GetCorpCache(tu.corpId, tu.agentTag, getType)
    client, req := tu.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := tu.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewTagUpdate(corpId, agentTag string) *tagUpdate {
    tu := &tagUpdate{wx.NewBaseWxCorp(), "", "", "", ""}
    tu.corpId = corpId
    tu.agentTag = agentTag
    tu.ReqContentType = project.HTTPContentTypeJSON
    tu.ReqMethod = fasthttp.MethodPost
    return tu
}
