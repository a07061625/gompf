/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 9:04
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

type unTaggingBatch struct {
    wx.BaseWxAccount
    appId      string
    tagId      int      // 标签ID
    openidList []string // 用户openid列表
}

func (utb *unTaggingBatch) SetTagId(tagId int) {
    if tagId > 0 {
        utb.tagId = tagId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "标签ID不合法", nil))
    }
}

func (utb *unTaggingBatch) SetOpenidList(openidList []string) {
    utb.openidList = make([]string, 0)
    for _, v := range openidList {
        match, _ := regexp.MatchString(project.RegexWxOpenid, v)
        if match {
            utb.openidList = append(utb.openidList, v)
        }
    }
}

func (utb *unTaggingBatch) checkData() {
    if utb.tagId <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "标签ID不能为空", nil))
    }
    if len(utb.openidList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid列表不能为空", nil))
    } else if len(utb.openidList) > 100 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid列表不能超过100个", nil))
    }
}

func (utb *unTaggingBatch) SendRequest() api.ApiResult {
    utb.checkData()

    reqData := make(map[string]interface{})
    reqData["tagid"] = utb.tagId
    reqData["openid_list"] = utb.openidList
    reqBody := mpf.JsonMarshal(reqData)
    utb.ReqUrl = "https://api.weixin.qq.com/cgi-bin/tags/members/batchuntagging?access_token=" + wx.NewUtilWx().GetSingleAccessToken(utb.appId)
    client, req := utb.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := utb.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewUnTaggingBatch(appId string) *unTaggingBatch {
    utb := &unTaggingBatch{wx.NewBaseWxAccount(), "", 0, make([]string, 0)}
    utb.appId = appId
    utb.ReqContentType = project.HttpContentTypeJson
    utb.ReqMethod = fasthttp.MethodPost
    return utb
}
