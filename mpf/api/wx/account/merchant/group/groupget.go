/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 16:25
 */
package group

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type groupGet struct {
    wx.BaseWxAccount
    appId   string
    groupId int // 分组ID
}

func (gg *groupGet) SetGroupId(groupId int) {
    if groupId > 0 {
        gg.groupId = groupId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分组ID不合法", nil))
    }
}

func (gg *groupGet) checkData() {
    if gg.groupId <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分组ID不能为空", nil))
    }
}

func (gg *groupGet) SendRequest() api.ApiResult {
    gg.checkData()

    reqData := make(map[string]interface{})
    reqData["group_id"] = gg.groupId
    reqBody := mpf.JsonMarshal(reqData)
    gg.ReqUrl = "https://api.weixin.qq.com/merchant/group/getbyid?access_token=" + wx.NewUtilWx().GetSingleAccessToken(gg.appId)
    client, req := gg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := gg.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewGroupGet(appId string) *groupGet {
    gg := &groupGet{wx.NewBaseWxAccount(), "", 0}
    gg.appId = appId
    gg.ReqContentType = project.HttpContentTypeJson
    gg.ReqMethod = fasthttp.MethodPost
    return gg
}
