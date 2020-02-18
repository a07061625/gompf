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

// 获取部门成员详情
type userList struct {
    wx.BaseWxCorp
    corpId         string
    agentTag       string
    departmentId   int // 部门id
    fetchChildFlag int // 匹配子部门标识 0:不匹配 1:匹配
}

func (ul *userList) SetDepartmentId(departmentId int) {
    if departmentId > 0 {
        ul.departmentId = departmentId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门id不合法", nil))
    }
}

func (ul *userList) SetFetchChildFlag(fetchChildFlag int) {
    if (fetchChildFlag == 0) || (fetchChildFlag == 1) {
        ul.ReqData["fetch_child"] = strconv.Itoa(fetchChildFlag)
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "匹配子部门标识不合法", nil))
    }
}

func (ul *userList) checkData() {
    if ul.departmentId <= 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门id不能为空", nil))
    }
    ul.ReqData["department_id"] = strconv.Itoa(ul.departmentId)
}

func (ul *userList) SendRequest(getType string) api.ApiResult {
    ul.checkData()

    ul.ReqData["access_token"] = wx.NewUtilWx().GetCorpCache(ul.corpId, ul.agentTag, getType)
    ul.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/user/list?" + mpf.HTTPCreateParams(ul.ReqData, "none", 1)
    client, req := ul.GetRequest()

    resp, result := ul.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewUserList(corpId, agentTag string) *userList {
    ul := &userList{wx.NewBaseWxCorp(), "", "", 0, 0}
    ul.corpId = corpId
    ul.agentTag = agentTag
    ul.ReqData["fetch_child"] = "0"
    return ul
}
