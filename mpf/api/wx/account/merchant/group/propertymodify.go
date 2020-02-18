/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 16:39
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

type propertyModify struct {
    wx.BaseWxAccount
    appId     string
    groupId   int    // 分组ID
    groupName string // 分组名称
}

func (pm *propertyModify) SetGroupId(groupId int) {
    if groupId > 0 {
        pm.groupId = groupId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分组ID不合法", nil))
    }
}

func (pm *propertyModify) SetGroupName(groupName string) {
    if len(groupName) > 0 {
        pm.groupName = groupName
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分组名称不合法", nil))
    }
}

func (pm *propertyModify) checkData() {
    if pm.groupId <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分组ID不能为空", nil))
    }
    if len(pm.groupName) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分组名称不能为空", nil))
    }
}

func (pm *propertyModify) SendRequest() api.ApiResult {
    pm.checkData()

    reqData := make(map[string]interface{})
    reqData["group_id"] = pm.groupId
    reqData["group_name"] = pm.groupName
    reqBody := mpf.JSONMarshal(reqData)
    pm.ReqUrl = "https://api.weixin.qq.com/merchant/group/propertymod?access_token=" + wx.NewUtilWx().GetSingleAccessToken(pm.appId)
    client, req := pm.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := pm.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewPropertyModify(appId string) *propertyModify {
    pm := &propertyModify{wx.NewBaseWxAccount(), "", 0, ""}
    pm.appId = appId
    pm.ReqContentType = project.HTTPContentTypeJSON
    pm.ReqMethod = fasthttp.MethodPost
    return pm
}
