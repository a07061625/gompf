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

// 更新部门
type departmentUpdate struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    name     string // 名称
    parentId int    // 父部门id
    orderNum int    // 排序值,数字越大越靠前
    id       int    // 部门id
}

func (du *departmentUpdate) SetName(name string) {
    if len(name) > 0 {
        trueName := []rune(name)
        du.name = string(trueName[:16])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "名称不合法", nil))
    }
}

func (du *departmentUpdate) SetParentId(parentId int) {
    if parentId >= 0 {
        du.parentId = parentId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "父部门id不合法", nil))
    }
}

func (du *departmentUpdate) SetOrderNum(orderNum int) {
    if orderNum >= 0 {
        du.orderNum = orderNum
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "排序值不合法", nil))
    }
}

func (du *departmentUpdate) SetId(id int) {
    if id > 0 {
        du.id = id
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门id不合法", nil))
    }
}

func (du *departmentUpdate) checkData() {
    if du.id <= 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门id不能为空", nil))
    }
}

func (du *departmentUpdate) SendRequest(getType string) api.ApiResult {
    du.checkData()
    reqData := make(map[string]interface{})
    reqData["id"] = du.id
    reqData["parentid"] = du.parentId
    reqData["order"] = du.orderNum
    if len(du.name) > 0 {
        reqData["name"] = du.name
    }
    reqBody := mpf.JsonMarshal(reqData)

    du.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/department/update?access_token=" + wx.NewUtilWx().GetCorpCache(du.corpId, du.agentTag, getType)
    client, req := du.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := du.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewDepartmentUpdate(corpId, agentTag string) *departmentUpdate {
    du := &departmentUpdate{wx.NewBaseWxCorp(), "", "", "", 0, 0, 0}
    du.corpId = corpId
    du.agentTag = agentTag
    du.ReqContentType = project.HTTPContentTypeJSON
    du.ReqMethod = fasthttp.MethodPost
    return du
}
