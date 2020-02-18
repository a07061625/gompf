/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 17:14
 */
package user

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 获取部门成员
type userSimpleList struct {
    wx.BaseWxCorp
    corpId         string
    agentTag       string
    departmentId   int // 部门id
    fetchChildFlag int // 匹配子部门标识 0:不匹配 1:匹配
}

func (usl *userSimpleList) SetDepartmentId(departmentId int) {
    if departmentId > 0 {
        usl.departmentId = departmentId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门id不合法", nil))
    }
}

func (usl *userSimpleList) SetFetchChildFlag(fetchChildFlag int) {
    if (fetchChildFlag == 0) || (fetchChildFlag == 1) {
        usl.ReqData["fetch_child"] = strconv.Itoa(fetchChildFlag)
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "匹配子部门标识不合法", nil))
    }
}

func (usl *userSimpleList) checkData() {
    if usl.departmentId <= 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门id不能为空", nil))
    }
    usl.ReqData["department_id"] = strconv.Itoa(usl.departmentId)
}

func (usl *userSimpleList) SendRequest(getType string) api.APIResult {
    usl.checkData()

    usl.ReqData["access_token"] = wx.NewUtilWx().GetCorpCache(usl.corpId, usl.agentTag, getType)
    usl.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist?" + mpf.HTTPCreateParams(usl.ReqData, "none", 1)
    client, req := usl.GetRequest()

    resp, resuslt := usl.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return resuslt
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        resuslt.Data = respData
    } else {
        resuslt.Code = errorcode.WxCorpRequestPost
        resuslt.Msg = respData["errmsg"].(string)
    }
    return resuslt
}

func NewUserSimpleList(corpId, agentTag string) *userSimpleList {
    usl := &userSimpleList{wx.NewBaseWxCorp(), "", "", 0, 0}
    usl.corpId = corpId
    usl.agentTag = agentTag
    usl.ReqData["fetch_child"] = "0"
    return usl
}
