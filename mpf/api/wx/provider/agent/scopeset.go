/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 23:01
 */
package agent

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

// 设置授权应用可见范围
type scopeSet struct {
    wx.BaseWxProvider
    accessToken    string   // 令牌,由查询注册状态接口返回
    agentId        int      // 应用ID
    allowUserList  []string // 应用成员可见范围
    allowPartyList []int    // 应用部门可见范围
    allowTagList   []int    // 应用标签可见范围
}

func (ss *scopeSet) SetAccessToken(accessToken string) {
    if len(accessToken) > 0 {
        ss.accessToken = accessToken
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "令牌不合法", nil))
    }
}

func (ss *scopeSet) SetAgentId(agentId int) {
    if agentId > 0 {
        ss.agentId = agentId
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "应用ID不合法", nil))
    }
}

func (ss *scopeSet) SetAllowUserList(allowUserList []string) {
    ss.allowUserList = make([]string, 0)
    for _, v := range allowUserList {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            ss.allowUserList = append(ss.allowUserList, v)
        }
    }
}

func (ss *scopeSet) SetAllowPartyList(allowPartyList []int) {
    ss.allowPartyList = make([]int, 0)
    for _, v := range allowPartyList {
        if v > 0 {
            ss.allowPartyList = append(ss.allowPartyList, v)
        }
    }
}

func (ss *scopeSet) SetAllowTagList(allowTagList []int) {
    ss.allowTagList = make([]int, 0)
    for _, v := range allowTagList {
        if v > 0 {
            ss.allowTagList = append(ss.allowTagList, v)
        }
    }
}

func (ss *scopeSet) checkData() {
    if len(ss.accessToken) == 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "令牌不能为空", nil))
    }
    if ss.agentId <= 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "应用ID不能为空", nil))
    }
}

func (ss *scopeSet) SendRequest() api.ApiResult {
    ss.checkData()

    reqData := make(map[string]interface{})
    reqData["agentid"] = ss.agentId
    if len(ss.allowUserList) > 0 {
        reqData["allow_user"] = ss.allowUserList
    }
    if len(ss.allowPartyList) > 0 {
        reqData["allow_party"] = ss.allowPartyList
    }
    if len(ss.allowTagList) > 0 {
        reqData["allow_tag"] = ss.allowTagList
    }
    reqBody := mpf.JsonMarshal(reqData)
    ss.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/agent/set_scope?access_token=" + ss.accessToken
    client, req := ss.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ss.SendInner(client, req, errorcode.WxProviderRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxProviderRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewScopeSet() *scopeSet {
    ss := &scopeSet{wx.NewBaseWxProvider(), "", 0, make([]string, 0), make([]int, 0), make([]int, 0)}
    ss.ReqContentType = project.HttpContentTypeJson
    ss.ReqMethod = fasthttp.MethodPost
    return ss
}
