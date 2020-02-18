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

// 获取部门列表
type departmentList struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    id       int // 部门id
}

func (dl *departmentList) SetId(id int) {
    if id > 0 {
        dl.ReqData["id"] = strconv.Itoa(id)
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门id不合法", nil))
    }
}

func (dl *departmentList) SendRequest(getType string) api.ApiResult {
    dl.ReqData["access_token"] = wx.NewUtilWx().GetCorpCache(dl.corpId, dl.agentTag, getType)
    dl.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/department/list?" + mpf.HTTPCreateParams(dl.ReqData, "none", 1)
    client, req := dl.GetRequest()

    resp, result := dl.SendInner(client, req, errorcode.WxCorpRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewDepartmentList(corpId, agentTag string) *departmentList {
    dl := &departmentList{wx.NewBaseWxCorp(), "", "", 0}
    dl.corpId = corpId
    dl.agentTag = agentTag
    return dl
}
