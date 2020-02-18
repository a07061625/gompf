/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 16:19
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

type groupDel struct {
    wx.BaseWxAccount
    appId   string
    groupId int // 分组ID
}

func (gd *groupDel) SetGroupId(groupId int) {
    if groupId > 0 {
        gd.groupId = groupId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分组ID不合法", nil))
    }
}

func (gd *groupDel) checkData() {
    if gd.groupId <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分组ID不能为空", nil))
    }
}

func (gd *groupDel) SendRequest() api.ApiResult {
    gd.checkData()

    reqData := make(map[string]interface{})
    reqData["group_id"] = gd.groupId
    reqBody := mpf.JSONMarshal(reqData)
    gd.ReqUrl = "https://api.weixin.qq.com/merchant/group/del?access_token=" + wx.NewUtilWx().GetSingleAccessToken(gd.appId)
    client, req := gd.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := gd.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewGroupDel(appId string) *groupDel {
    gd := &groupDel{wx.NewBaseWxAccount(), "", 0}
    gd.appId = appId
    gd.ReqContentType = project.HTTPContentTypeJSON
    gd.ReqMethod = fasthttp.MethodPost
    return gd
}
