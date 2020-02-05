/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 9:41
 */
package department

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 删除部门
type departmentDelete struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    id       int // 部门id
}

func (dd *departmentDelete) SetId(id int) {
    if id > 0 {
        dd.id = id
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门id不合法", nil))
    }
}

func (dd *departmentDelete) checkData() {
    if dd.id <= 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门id不能为空", nil))
    }
    dd.ReqData["id"] = strconv.Itoa(dd.id)
}

func (dd *departmentDelete) SendRequest(getType string) api.ApiResult {
    dd.checkData()

    dd.ReqData["access_token"] = wx.NewUtilWx().GetCorpCache(dd.corpId, dd.agentTag, getType)
    dd.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/department/delete?" + mpf.HttpCreateParams(dd.ReqData, "none", 1)
    client, req := dd.GetRequest()

    resp, result := dd.SendInner(client, req, errorcode.WxCorpRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewDepartmentDelete(corpId, agentTag string) *departmentDelete {
    dd := &departmentDelete{wx.NewBaseWxCorp(), "", "", 0}
    dd.corpId = corpId
    dd.agentTag = agentTag
    return dd
}
