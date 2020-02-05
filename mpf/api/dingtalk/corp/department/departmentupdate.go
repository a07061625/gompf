package department

import (
    "regexp"
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 更新部门
type departmentUpdate struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    id       int // 部门id
}

func (du *departmentUpdate) SetId(id int) {
    if id > 0 {
        du.id = id
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不合法", nil))
    }
}

func (du *departmentUpdate) SetName(name string) {
    if len(name) > 0 {
        trueName := []rune(name)
        du.ExtendData["name"] = string(trueName[:32])
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "名称不合法", nil))
    }
}

func (du *departmentUpdate) SetLang(lang string) {
    if (lang == "zh_CN") || (lang == "en_US") {
        du.ExtendData["lang"] = lang
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "语言不合法", nil))
    }
}

func (du *departmentUpdate) SetParentId(parentId int) {
    if parentId > 0 {
        du.ExtendData["parentid"] = parentId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "父部门id不合法", nil))
    }
}

func (du *departmentUpdate) SetOrder(order int) {
    if order >= 0 {
        du.ExtendData["order"] = order
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "排序值不合法", nil))
    }
}

func (du *departmentUpdate) SetCreateDeptGroup(createDeptGroup bool) {
    du.ExtendData["createDeptGroup"] = createDeptGroup
}

func (du *departmentUpdate) SetGroupContainSubDept(groupContainSubDept bool) {
    du.ExtendData["groupContainSubDept"] = groupContainSubDept
}

func (du *departmentUpdate) SetGroupContainOuterDept(groupContainOuterDept bool) {
    du.ExtendData["groupContainOuterDept"] = groupContainOuterDept
}

func (du *departmentUpdate) SetGroupContainHiddenDept(groupContainHiddenDept bool) {
    du.ExtendData["groupContainHiddenDept"] = groupContainHiddenDept
}

func (du *departmentUpdate) SetAutoAddUser(autoAddUser bool) {
    du.ExtendData["autoAddUser"] = autoAddUser
}

func (du *departmentUpdate) SetDeptManagerUserList(deptManagerUserList []string) {
    userList := make([]string, 0)
    for _, v := range deptManagerUserList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            userList = append(userList, v)
        }
    }

    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门主管列表不能为空", nil))
    } else if len(userList) > 200 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门主管列表不能超过200个", nil))
    }
    du.ExtendData["deptManagerUseridList"] = strings.Join(userList, "|")
}

func (du *departmentUpdate) SetDeptHiding(deptHiding bool) {
    du.ExtendData["deptHiding"] = deptHiding
}

func (du *departmentUpdate) SetDeptPermits(deptPermits []int) {
    departmentList := make([]string, 0)
    for _, v := range deptPermits {
        if v > 0 {
            departmentList = append(departmentList, strconv.Itoa(v))
        }
    }

    if len(departmentList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门列表不能为空", nil))
    } else if len(departmentList) > 200 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门列表不能超过200个", nil))
    }
    du.ExtendData["deptPermits"] = strings.Join(departmentList, "|")
}

func (du *departmentUpdate) SetUserPermits(userPermits []string) {
    userList := make([]string, 0)
    for _, v := range userPermits {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            userList = append(userList, v)
        }
    }

    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "人员列表不能为空", nil))
    } else if len(userList) > 200 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "人员列表不能超过200个", nil))
    }
    du.ExtendData["userPermits"] = strings.Join(userList, "|")
}

func (du *departmentUpdate) SetOuterDept(outerDept bool) {
    du.ExtendData["outerDept"] = outerDept
}

func (du *departmentUpdate) SetOuterPermitDepts(outerPermitDepts []int) {
    departmentList := make([]string, 0)
    for _, v := range outerPermitDepts {
        if v > 0 {
            departmentList = append(departmentList, strconv.Itoa(v))
        }
    }

    if len(departmentList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "额外部门列表不能为空", nil))
    } else if len(departmentList) > 200 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "额外部门列表不能超过200个", nil))
    }
    du.ExtendData["outerPermitDepts"] = strings.Join(departmentList, "|")
}

func (du *departmentUpdate) SetOuterPermitUsers(outerPermitUsers []string) {
    userList := make([]string, 0)
    for _, v := range outerPermitUsers {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            userList = append(userList, v)
        }
    }

    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "额外人员列表不能为空", nil))
    } else if len(userList) > 200 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "额外人员列表不能超过200个", nil))
    }
    du.ExtendData["outerPermitUsers"] = strings.Join(userList, "|")
}

func (du *departmentUpdate) SetOuterDeptOnlySelf(outerDeptOnlySelf bool) {
    du.ExtendData["outerDeptOnlySelf"] = outerDeptOnlySelf
}

func (du *departmentUpdate) SetOrgDeptOwner(orgDeptOwner string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, orgDeptOwner)
    if match {
        du.ExtendData["orgDeptOwner"] = orgDeptOwner
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "群主不合法", nil))
    }
}

func (du *departmentUpdate) SetSourceIdentifier(sourceIdentifier string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, sourceIdentifier)
    if match {
        du.ExtendData["sourceIdentifier"] = sourceIdentifier
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门标识不合法", nil))
    }
}

func (du *departmentUpdate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if du.id <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不能为空", nil))
    }
    du.ExtendData["id"] = du.id

    du.ReqUrl = dingtalk.UrlService + "/department/update?access_token=" + dingtalk.NewUtil().GetAccessToken(du.corpId, du.agentTag, du.atType)

    reqBody := mpf.JsonMarshal(du.ExtendData)
    client, req := du.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewDepartmentUpdate(corpId, agentTag, atType string) *departmentUpdate {
    du := &departmentUpdate{dingtalk.NewCorp(), "", "", "", 0}
    du.corpId = corpId
    du.agentTag = agentTag
    du.atType = atType
    du.ReqContentType = project.HttpContentTypeJson
    du.ReqMethod = fasthttp.MethodPost
    return du
}
