/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 9:41
 */
package department

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建部门
type departmentCreate struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    name     string // 名称
    parentId int    // 父部门id
    orderNum int    // 排序值,数字越大越靠前
    id       int    // 部门id
}

func (dc *departmentCreate) SetName(name string) {
    if len(name) > 0 {
        trueName := []rune(name)
        dc.name = string(trueName[:16])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "名称不合法", nil))
    }
}

func (dc *departmentCreate) SetParentId(parentId int) {
    if parentId >= 0 {
        dc.parentId = parentId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "父部门id不合法", nil))
    }
}

func (dc *departmentCreate) SetOrderNum(orderNum int) {
    if orderNum >= 0 {
        dc.orderNum = orderNum
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "排序值不合法", nil))
    }
}

func (dc *departmentCreate) SetId(id int) {
    if id > 0 {
        dc.id = id
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门id不合法", nil))
    }
}

func (dc *departmentCreate) checkData() {
    if len(dc.name) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "名称不能为空", nil))
    }
}

func (dc *departmentCreate) SendRequest(getType string) api.ApiResult {
    dc.checkData()
    reqData := make(map[string]interface{})
    reqData["name"] = dc.name
    reqData["parentid"] = dc.parentId
    reqData["order"] = dc.orderNum
    if dc.id > 0 {
        reqData["id"] = dc.id
    }
    reqBody := mpf.JsonMarshal(reqData)

    dc.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/department/create?access_token=" + wx.NewUtilWx().GetCorpCache(dc.corpId, dc.agentTag, getType)
    client, req := dc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := dc.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewDepartmentCreate(corpId, agentTag string) *departmentCreate {
    dc := &departmentCreate{wx.NewBaseWxCorp(), "", "", "", 0, 0, 0}
    dc.corpId = corpId
    dc.agentTag = agentTag
    dc.ReqContentType = project.HttpContentTypeJson
    dc.ReqMethod = fasthttp.MethodPost
    return dc
}
