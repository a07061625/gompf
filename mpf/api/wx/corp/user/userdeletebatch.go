/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 17:02
 */
package user

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 批量删除成员
type userDeleteBatch struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    userList []string // 用户列表
}

func (udb *userDeleteBatch) SetUserList(userList []string) {
    udb.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            udb.userList = append(udb.userList, v)
        }
    }
    if len(udb.userList) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户列表不能为空", nil))
    } else if len(udb.userList) > 200 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户列表不能超过200个", nil))
    }
}

func (udb *userDeleteBatch) checkData() {
    if len(udb.userList) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户列表不能为空", nil))
    }
}

func (udb *userDeleteBatch) SendRequest(getType string) api.APIResult {
    udb.checkData()

    reqData := make(map[string]interface{})
    reqData["useridlist"] = udb.userList
    reqBody := mpf.JSONMarshal(reqData)
    udb.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/user/batchdelete?access_token=" + wx.NewUtilWx().GetCorpCache(udb.corpId, udb.agentTag, getType)
    client, req := udb.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := udb.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewUserDeleteBatch(corpId, agentTag string) *userDeleteBatch {
    udb := &userDeleteBatch{wx.NewBaseWxCorp(), "", "", make([]string, 0)}
    udb.corpId = corpId
    udb.agentTag = agentTag
    udb.ReqContentType = project.HTTPContentTypeJSON
    udb.ReqMethod = fasthttp.MethodPost
    return udb
}
