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

// 取消拉黑用户
type unBlackBatch struct {
    wx.BaseWxAccount
    appId      string
    openidList []string // 用户openid列表
}

func (ubb *unBlackBatch) SetOpenidList(openidList []string) {
    ubb.openidList = make([]string, 0)
    for _, v := range openidList {
        match, _ := regexp.MatchString(project.RegexWxOpenid, v)
        if match {
            ubb.openidList = append(ubb.openidList, v)
        }
    }
}

func (ubb *unBlackBatch) checkData() {
    if len(ubb.openidList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid列表不能为空", nil))
    } else if len(ubb.openidList) > 100 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid列表不能超过100个", nil))
    }
}

func (ubb *unBlackBatch) SendRequest() api.ApiResult {
    ubb.checkData()

    reqData := make(map[string]interface{})
    reqData["openid_list"] = ubb.openidList
    reqBody := mpf.JsonMarshal(reqData)
    ubb.ReqUrl = "https://api.weixin.qq.com/cgi-bin/tags/members/batchunblacklist?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ubb.appId)
    client, req := ubb.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ubb.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewUnBlackBatch(appId string) *unBlackBatch {
    ubb := &unBlackBatch{wx.NewBaseWxAccount(), "", make([]string, 0)}
    ubb.appId = appId
    ubb.ReqContentType = project.HttpContentTypeJson
    ubb.ReqMethod = fasthttp.MethodPost
    return ubb
}
