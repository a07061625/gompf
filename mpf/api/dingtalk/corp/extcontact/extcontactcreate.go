package extcontact

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

// 添加外部联系人
type extContactCreate struct {
    dingtalk.BaseCorp
    corpId         string
    agentTag       string
    labelList      []int  // 标签列表
    followerUserId string // 负责人用户id
    name           string // 名称
    mobile         string // 手机号
}

func (ecc *extContactCreate) SetLabelList(labelList []int) {
    ecc.labelList = make([]int, 0)
    for _, v := range labelList {
        if v > 0 {
            ecc.labelList = append(ecc.labelList, v)
        }
    }

    if len(ecc.labelList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "标签列表不能为空", nil))
    }
}

func (ecc *extContactCreate) SetFollowerUserId(followerUserId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, followerUserId)
    if match {
        ecc.followerUserId = followerUserId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "负责人用户id不合法", nil))
    }
}

func (ecc *extContactCreate) SetName(name string) {
    if len(name) > 0 {
        ecc.name = name
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "名称不合法", nil))
    }
}

func (ecc *extContactCreate) SetMobile(mobile string) {
    match, _ := regexp.MatchString(project.RegexPhone, mobile)
    if match {
        ecc.mobile = mobile
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "手机号不合法", nil))
    }
}

func (ecc *extContactCreate) SetTitle(title string) {
    ecc.ExtendData["title"] = strings.TrimSpace(title)
}

func (ecc *extContactCreate) SetShareDeptIds(shareDeptIds []int) {
    deptList := make([]int, 0)
    for _, v := range shareDeptIds {
        if v > 0 {
            deptList = append(deptList, v)
        }
    }
    if len(deptList) > 0 {
        ecc.ExtendData["share_dept_ids"] = deptList
    }
}

func (ecc *extContactCreate) SetAddress(address string) {
    ecc.ExtendData["address"] = strings.TrimSpace(address)
}

func (ecc *extContactCreate) SetRemark(remark string) {
    ecc.ExtendData["remark"] = strings.TrimSpace(remark)
}

func (ecc *extContactCreate) SetCompanyName(companyName string) {
    ecc.ExtendData["company_name"] = strings.TrimSpace(companyName)
}

func (ecc *extContactCreate) SetShareUserIds(shareUserIds []string) {
    userList := make([]string, 0)
    for _, v := range shareUserIds {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            userList = append(userList, v)
        }
    }
    if len(userList) > 0 {
        ecc.ExtendData["share_user_ids"] = userList
    }
}

func (ecc *extContactCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ecc.labelList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "标签列表不能为空", nil))
    }
    if len(ecc.followerUserId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "负责人用户id不能为空", nil))
    }
    if len(ecc.name) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "名称不能为空", nil))
    }
    if len(ecc.mobile) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "手机号不能为空", nil))
    }
    ecc.ExtendData["label_ids"] = ecc.labelList
    ecc.ExtendData["follower_user_id"] = ecc.followerUserId
    ecc.ExtendData["name"] = ecc.name
    ecc.ExtendData["mobile"] = ecc.mobile

    ecc.ReqUrl = dingtalk.UrlService + "/topapi/extcontact/create?access_token=" + dingtalk.NewUtil().GetCorpAccessToken(ecc.corpId, ecc.agentTag)

    reqData := make(map[string]interface{})
    reqData["contact"] = ecc.ExtendData
    reqBody := mpf.JsonMarshal(reqData)
    client, req := ecc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewExtContactCreate(corpId, agentTag string) *extContactCreate {
    ecc := &extContactCreate{dingtalk.NewCorp(), "", "", make([]int, 0), "", "", ""}
    ecc.corpId = corpId
    ecc.agentTag = agentTag
    ecc.ExtendData["state_code"] = "86"
    ecc.ReqContentType = project.HTTPContentTypeJSON
    ecc.ReqMethod = fasthttp.MethodPost
    return ecc
}
