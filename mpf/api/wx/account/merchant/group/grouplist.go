/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 16:22
 */
package group

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

type groupList struct {
    wx.BaseWxAccount
    appId string
}

func (gl *groupList) SendRequest() api.ApiResult {
    gl.ReqUrl = "https://api.weixin.qq.com/merchant/group/getall?access_token=" + wx.NewUtilWx().GetSingleAccessToken(gl.appId)
    client, req := gl.GetRequest()

    resp, result := gl.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewGroupList(appId string) *groupList {
    gl := &groupList{wx.NewBaseWxAccount(), ""}
    gl.appId = appId
    return gl
}
