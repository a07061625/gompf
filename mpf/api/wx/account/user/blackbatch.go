/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 23:44
 */
package user

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

// 拉黑用户
type blackBatch struct {
    wx.BaseWxAccount
    appId      string
    openidList []string // 用户openid列表
}

func (bb *blackBatch) SetOpenidList(openidList []string) {
    bb.openidList = make([]string, 0)
    for _, v := range openidList {
        match, _ := regexp.MatchString(project.RegexWxOpenid, v)
        if match {
            bb.openidList = append(bb.openidList, v)
        }
    }
}

func (bb *blackBatch) checkData() {
    if len(bb.openidList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid列表不能为空", nil))
    } else if len(bb.openidList) > 100 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid列表不能超过100个", nil))
    }
}

func (bb *blackBatch) SendRequest() api.ApiResult {
    bb.checkData()

    reqData := make(map[string]interface{})
    reqData["openid_list"] = bb.openidList
    reqBody := mpf.JsonMarshal(reqData)
    bb.ReqUrl = "https://api.weixin.qq.com/cgi-bin/tags/members/batchblacklist?access_token=" + wx.NewUtilWx().GetSingleAccessToken(bb.appId)
    client, req := bb.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := bb.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewBlackBatch(appId string) *blackBatch {
    bb := &blackBatch{wx.NewBaseWxAccount(), "", make([]string, 0)}
    bb.appId = appId
    bb.ReqContentType = project.HTTPContentTypeJSON
    bb.ReqMethod = fasthttp.MethodPost
    return bb
}
