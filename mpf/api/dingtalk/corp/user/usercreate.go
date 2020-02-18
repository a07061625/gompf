/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 1:05
 */
package user

import (
    "regexp"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建用户
type userCreate struct {
    dingtalk.BaseCorp
    corpId         string
    agentTag       string
    atType         string
    name           string // 名称
    departmentList []int  // 部门列表
    mobile         string // 手机号码
    jobNumber      string // 工号
}

func (uc *userCreate) SetName(name string) {
    if len(name) > 0 {
        trueName := []rune(name)
        uc.name = string(trueName[:32])
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "名称不合法", nil))
    }
}

func (uc *userCreate) SetDepartmentList(departmentList []int) {
    if len(departmentList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门列表不能为空", nil))
    }

    uc.departmentList = make([]int, 0)
    for _, v := range departmentList {
        if v > 0 {
            uc.departmentList = append(uc.departmentList, v)
        }
    }
}

func (uc *userCreate) SetMobile(mobile string) {
    match, _ := regexp.MatchString(project.RegexPhone, mobile)
    if match {
        uc.mobile = mobile
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "手机号码不合法", nil))
    }
}

func (uc *userCreate) SetJobNumber(jobNumber string) {
    if (len(jobNumber) > 0) && (len(jobNumber) <= 64) {
        uc.jobNumber = jobNumber
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "工号不合法", nil))
    }
}

func (uc *userCreate) SetUserId(userId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,64}$`, userId)
    if match {
        uc.ExtendData["userid"] = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不合法", nil))
    }
}

func (uc *userCreate) SetDepartmentOrder(departmentOrder map[int]int) {
    if len(departmentOrder) > 0 {
        uc.ExtendData["orderInDepts"] = mpf.JSONMarshal(departmentOrder)
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门排序不合法", nil))
    }
}

func (uc *userCreate) SetPosition(position string) {
    if len(position) > 0 {
        truePosition := []rune(position)
        uc.ExtendData["position"] = string(truePosition[:32])
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "职位信息不合法", nil))
    }
}

func (uc *userCreate) SetTel(tel string) {
    if (len(tel) > 0) && (len(tel) <= 50) {
        uc.ExtendData["tel"] = tel
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分机号不合法", nil))
    }
}

func (uc *userCreate) SetWorkPlace(workPlace string) {
    if len(workPlace) > 0 {
        trueWorkPlace := []rune(workPlace)
        uc.ExtendData["workPlace"] = string(trueWorkPlace[:25])
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "办公地点不合法", nil))
    }
}

func (uc *userCreate) SetRemark(remark string) {
    if len(remark) > 0 {
        trueRemark := []rune(remark)
        uc.ExtendData["remark"] = string(trueRemark[:500])
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "备注不合法", nil))
    }
}

func (uc *userCreate) SetEmail(email string) {
    if (len(email) > 0) && (len(email) <= 64) {
        uc.ExtendData["email"] = email
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "邮箱不合法", nil))
    }
}

func (uc *userCreate) SetOrgEmail(orgEmail string) {
    if len(orgEmail) > 0 {
        uc.ExtendData["orgEmail"] = orgEmail
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "企业邮箱不合法", nil))
    }
}

func (uc *userCreate) SetIsHide(isHide bool) {
    uc.ExtendData["isHide"] = isHide
}

func (uc *userCreate) SetIsSenior(isSenior bool) {
    uc.ExtendData["isSenior"] = isSenior
}

func (uc *userCreate) SetExtAttr(extAttr map[string]interface{}) {
    if len(extAttr) > 0 {
        uc.ExtendData["extattr"] = extAttr
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "扩展属性不合法", nil))
    }
}

func (uc *userCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(uc.name) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "名称不能为空", nil))
    }
    if len(uc.departmentList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门列表不能为空", nil))
    }
    if len(uc.mobile) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "手机号码不能为空", nil))
    }
    if len(uc.jobNumber) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "工号不能为空", nil))
    }
    uc.ExtendData["name"] = uc.name
    uc.ExtendData["department"] = uc.departmentList
    uc.ExtendData["mobile"] = uc.mobile
    uc.ExtendData["jobnumber"] = uc.jobNumber

    uc.ReqUrl = dingtalk.UrlService + "/user/create?access_token=" + dingtalk.NewUtil().GetAccessToken(uc.corpId, uc.agentTag, uc.atType)

    reqBody := mpf.JSONMarshal(uc.ExtendData)
    client, req := uc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewUserCreate(corpId, agentTag, atType string) *userCreate {
    uc := &userCreate{dingtalk.NewCorp(), "", "", "", "", make([]int, 0), "", ""}
    uc.corpId = corpId
    uc.agentTag = agentTag
    uc.atType = atType
    uc.ExtendData["userid"] = mpf.ToolCreateNonceStr(8, "numlower") + strconv.FormatInt(time.Now().Unix(), 10)
    uc.ExtendData["isHide"] = false
    uc.ExtendData["isSenior"] = false
    uc.ReqContentType = project.HTTPContentTypeJSON
    uc.ReqMethod = fasthttp.MethodPost
    return uc
}
