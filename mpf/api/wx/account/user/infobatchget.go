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

type infoBatchGet struct {
    wx.BaseWxAccount
    appId      string
    openidList []string // 用户openid列表
}

func (ibg *infoBatchGet) SetOpenidList(openidList []string) {
    ibg.openidList = make([]string, 0)
    for _, v := range openidList {
        match, _ := regexp.MatchString(project.RegexWxOpenid, v)
        if match {
            ibg.openidList = append(ibg.openidList, v)
        }
    }
}

func (ibg *infoBatchGet) checkData() {
    if len(ibg.openidList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid列表不能为空", nil))
    } else if len(ibg.openidList) > 100 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid列表不能超过100个", nil))
    }
}

func (ibg *infoBatchGet) SendRequest() api.ApiResult {
    ibg.checkData()

    userList := make([]map[string]string, 0)
    for _, v := range ibg.openidList {
        userInfo := make(map[string]string)
        userInfo["openid"] = v
        userInfo["lang"] = "zh-CN"
        userList = append(userList, userInfo)
    }
    reqData := make(map[string]interface{})
    reqData["user_list"] = userList
    reqBody := mpf.JSONMarshal(reqData)
    ibg.ReqUrl = "https://api.weixin.qq.com/cgi-bin/user/info/batchget?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ibg.appId)
    client, req := ibg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ibg.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["user_info_list"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewInfoBatchGet(appId string) *infoBatchGet {
    ibg := &infoBatchGet{wx.NewBaseWxAccount(), "", make([]string, 0)}
    ibg.appId = appId
    ibg.ReqContentType = project.HTTPContentTypeJSON
    ibg.ReqMethod = fasthttp.MethodPost
    return ibg
}
