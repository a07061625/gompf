/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 1:05
 */
package user

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 更新用户
type userUpdate struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userId   string // 用户id
}

func (uu *userUpdate) SetUserId(userId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,64}$`, userId)
    if match {
        uu.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不合法", nil))
    }
}

func (uu *userUpdate) SetLang(lang string) {
    if (lang == "zh_CN") || (lang == "en_US") {
        uu.ExtendData["lang"] = lang
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "语言不合法", nil))
    }
}

func (uu *userUpdate) SetName(name string) {
    if len(name) > 0 {
        trueName := []rune(name)
        uu.ExtendData["name"] = string(trueName[:32])
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "名称不合法", nil))
    }
}

func (uu *userUpdate) SetDepartmentList(departmentList []int) {
    if len(departmentList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门列表不能为空", nil))
    }

    departments := make([]int, 0)
    for _, v := range departmentList {
        if v > 0 {
            departments = append(departments, v)
        }
    }
    if len(departments) > 0 {
        uu.ExtendData["department"] = departments
    }
}

func (uu *userUpdate) SetMobile(mobile string) {
    match, _ := regexp.MatchString(project.RegexPhone, mobile)
    if match {
        uu.ExtendData["mobile"] = mobile
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "手机号码不合法", nil))
    }
}

func (uu *userUpdate) SetJobNumber(jobNumber string) {
    if (len(jobNumber) > 0) && (len(jobNumber) <= 64) {
        uu.ExtendData["jobnumber"] = jobNumber
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "工号不合法", nil))
    }
}

func (uu *userUpdate) SetDepartmentOrder(departmentOrder map[int]int) {
    if len(departmentOrder) > 0 {
        uu.ExtendData["orderInDepts"] = mpf.JSONMarshal(departmentOrder)
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门排序不合法", nil))
    }
}

func (uu *userUpdate) SetPosition(position string) {
    if len(position) > 0 {
        truePosition := []rune(position)
        uu.ExtendData["position"] = string(truePosition[:32])
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "职位信息不合法", nil))
    }
}

func (uu *userUpdate) SetTel(tel string) {
    if (len(tel) > 0) && (len(tel) <= 50) {
        uu.ExtendData["tel"] = tel
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分机号不合法", nil))
    }
}

func (uu *userUpdate) SetWorkPlace(workPlace string) {
    if len(workPlace) > 0 {
        trueWorkPlace := []rune(workPlace)
        uu.ExtendData["workPlace"] = string(trueWorkPlace[:25])
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "办公地点不合法", nil))
    }
}

func (uu *userUpdate) SetRemark(remark string) {
    if len(remark) > 0 {
        trueRemark := []rune(remark)
        uu.ExtendData["remark"] = string(trueRemark[:500])
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "备注不合法", nil))
    }
}

func (uu *userUpdate) SetEmail(email string) {
    if (len(email) > 0) && (len(email) <= 64) {
        uu.ExtendData["email"] = email
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "邮箱不合法", nil))
    }
}

func (uu *userUpdate) SetOrgEmail(orgEmail string) {
    if len(orgEmail) > 0 {
        uu.ExtendData["orgEmail"] = orgEmail
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "企业邮箱不合法", nil))
    }
}

func (uu *userUpdate) SetIsHide(isHide bool) {
    uu.ExtendData["isHide"] = isHide
}

func (uu *userUpdate) SetIsSenior(isSenior bool) {
    uu.ExtendData["isSenior"] = isSenior
}

func (uu *userUpdate) SetExtAttr(extAttr map[string]interface{}) {
    if len(extAttr) > 0 {
        uu.ExtendData["extattr"] = extAttr
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "扩展属性不合法", nil))
    }
}

func (uu *userUpdate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(uu.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不能为空", nil))
    }
    uu.ExtendData["userid"] = uu.userId

    uu.ReqURI = dingtalk.UrlService + "/user/update?access_token=" + dingtalk.NewUtil().GetAccessToken(uu.corpId, uu.agentTag, uu.atType)

    reqBody := mpf.JSONMarshal(uu.ExtendData)
    client, req := uu.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewUserUpdate(corpId, agentTag, atType string) *userUpdate {
    uu := &userUpdate{dingtalk.NewCorp(), "", "", "", ""}
    uu.corpId = corpId
    uu.agentTag = agentTag
    uu.atType = atType
    uu.ReqContentType = project.HTTPContentTypeJSON
    uu.ReqMethod = fasthttp.MethodPost
    return uu
}
