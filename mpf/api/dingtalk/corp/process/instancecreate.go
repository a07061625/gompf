package process

import (
    "regexp"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 发起审批实例
type instanceCreate struct {
    dingtalk.BaseCorp
    corpId        string
    agentTag      string
    atType        string
    processCode   string                   // 审批码
    userId        string                   // 发起人用户ID
    departId      int                      // 发起人部门ID
    approverList  []string                 // 审批人列表
    formValueList []map[string]interface{} // 表单参数列表
}

func (ic *instanceCreate) SetProcessCode(processCode string) {
    if len(processCode) > 0 {
        ic.processCode = processCode
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "审批码不合法", nil))
    }
}

func (ic *instanceCreate) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        ic.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "发起人用户ID不合法", nil))
    }
}

func (ic *instanceCreate) SetDepartId(departId int) {
    if (departId == -1) || (departId > 0) {
        ic.departId = departId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "发起人部门ID不合法", nil))
    }
}

func (ic *instanceCreate) SetApproverList(approverList []string) {
    if len(approverList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "审批人列表不能为空", nil))
    } else if len(approverList) > 20 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "审批人不能超过20个", nil))
    }

    approvers := make([]string, 0)
    for _, v := range approverList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            approvers = append(approvers, v)
        }
    }
    ic.approverList = approvers
}

func (ic *instanceCreate) SetFormValueList(formValueList []map[string]interface{}) {
    if len(formValueList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "表单参数不能为空", nil))
    } else if len(formValueList) > 20 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "表单参数不能超过20个", nil))
    }

    ic.formValueList = formValueList
}

func (ic *instanceCreate) SetCcList(ccList []string) {
    if len(ccList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "抄送人列表不能为空", nil))
    } else if len(ccList) > 20 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "抄送人不能超过20个", nil))
    }

    userList := make([]string, 0)
    for _, v := range ccList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            userList = append(userList, v)
        }
    }
    ic.ExtendData["cc_list"] = strings.Join(userList, ",")
}

func (ic *instanceCreate) SetCcPosition(ccPosition string) {
    if len(ccPosition) > 0 {
        ic.ExtendData["cc_position"] = ccPosition
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "抄送时间不合法", nil))
    }
}

func (ic *instanceCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ic.processCode) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "审批码不能为空", nil))
    }
    if len(ic.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "发起人用户ID不能为空", nil))
    }
    if ic.departId == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "发起人部门ID不能为空", nil))
    }
    if len(ic.approverList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "审批人列表不能为空", nil))
    }
    if len(ic.formValueList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "表单参数列表不能为空", nil))
    }
    ic.ExtendData["process_code"] = ic.processCode
    ic.ExtendData["originator_user_id"] = ic.userId
    ic.ExtendData["dept_id"] = ic.departId
    ic.ExtendData["approvers"] = strings.Join(ic.approverList, ",")
    ic.ExtendData["form_component_values"] = ic.formValueList
    if ic.atType == dingtalk.AccessTokenTypeProviderAuthorize {
        ic.ExtendData["agent_id"] = dingtalk.NewConfig().GetProvider().GetSuiteId()
    }

    ic.ReqUrl = dingtalk.UrlService + "/topapi/processinstance/create?access_token=" + dingtalk.NewUtil().GetAccessToken(ic.corpId, ic.agentTag, ic.atType)

    reqBody := mpf.JSONMarshal(ic.ExtendData)
    client, req := ic.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewInstanceCreate(corpId, agentTag, atType string) *instanceCreate {
    ic := &instanceCreate{dingtalk.NewCorp(), "", "", "", "", "", 0, make([]string, 0), make([]map[string]interface{}, 0)}
    ic.corpId = corpId
    ic.agentTag = agentTag
    ic.atType = atType
    ic.ReqContentType = project.HTTPContentTypeJSON
    ic.ReqMethod = fasthttp.MethodPost
    return ic
}
