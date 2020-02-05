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

// 创建部门
type departmentCreate struct {
    dingtalk.BaseCorp
    corpId           string
    agentTag         string
    atType           string
    name             string // 名称
    parentId         int    // 父部门id
    sourceIdentifier string // 部门标识
}

func (dc *departmentCreate) SetName(name string) {
    if len(name) > 0 {
        trueName := []rune(name)
        dc.name = string(trueName[:32])
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "名称不合法", nil))
    }
}

func (dc *departmentCreate) SetParentId(parentId int) {
    if parentId > 0 {
        dc.parentId = parentId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "父部门id不合法", nil))
    }
}

func (dc *departmentCreate) SetSourceIdentifier(sourceIdentifier string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, sourceIdentifier)
    if match {
        dc.sourceIdentifier = sourceIdentifier
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门标识不合法", nil))
    }
}

func (dc *departmentCreate) SetOrder(order int) {
    if order >= 0 {
        dc.ExtendData["order"] = order
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "排序值不合法", nil))
    }
}

func (dc *departmentCreate) SetDeptHiding(deptHiding bool) {
    dc.ExtendData["deptHiding"] = deptHiding
}

func (dc *departmentCreate) SetCreateDeptGroup(createDeptGroup bool) {
    dc.ExtendData["createDeptGroup"] = createDeptGroup
}

func (dc *departmentCreate) SetDeptPermits(deptPermits []int) {
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
    dc.ExtendData["deptPermits"] = strings.Join(departmentList, "|")
}

func (dc *departmentCreate) SetUserPermits(userPermits []string) {
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
    dc.ExtendData["userPermits"] = strings.Join(userList, "|")
}

func (dc *departmentCreate) SetOuterPermitDepts(outerPermitDepts []int) {
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
    dc.ExtendData["outerPermitDepts"] = strings.Join(departmentList, "|")
}

func (dc *departmentCreate) SetOuterPermitUsers(outerPermitUsers []string) {
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
    dc.ExtendData["outerPermitUsers"] = strings.Join(userList, "|")
}

func (dc *departmentCreate) SetOuterDept(outerDept bool) {
    dc.ExtendData["outerDept"] = outerDept
}

func (dc *departmentCreate) SetOuterDeptOnlySelf(outerDeptOnlySelf bool) {
    dc.ExtendData["outerDeptOnlySelf"] = outerDeptOnlySelf
}

func (dc *departmentCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dc.name) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "名称不能为空", nil))
    }
    if dc.parentId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "父部门id不能为空", nil))
    }
    if len(dc.sourceIdentifier) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门标识不能为空	", nil))
    }
    dc.ExtendData["name"] = dc.name
    dc.ExtendData["parentid"] = dc.parentId
    dc.ExtendData["sourceIdentifier"] = dc.sourceIdentifier

    dc.ReqUrl = dingtalk.UrlService + "/department/create?access_token=" + dingtalk.NewUtil().GetAccessToken(dc.corpId, dc.agentTag, dc.atType)

    reqBody := mpf.JsonMarshal(dc.ExtendData)
    client, req := dc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewDepartmentCreate(corpId, agentTag, atType string) *departmentCreate {
    dc := &departmentCreate{dingtalk.NewCorp(), "", "", "", "", 0, ""}
    dc.corpId = corpId
    dc.agentTag = agentTag
    dc.atType = atType
    dc.ExtendData["order"] = 0
    dc.ExtendData["createDeptGroup"] = false
    dc.ExtendData["deptHiding"] = false
    dc.ExtendData["outerDept"] = false
    dc.ExtendData["outerDeptOnlySelf"] = false
    dc.ReqContentType = project.HttpContentTypeJson
    dc.ReqMethod = fasthttp.MethodPost
    return dc
}
