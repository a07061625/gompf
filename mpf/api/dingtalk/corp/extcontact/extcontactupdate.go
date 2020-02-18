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

// 更新外部联系人
type extContactUpdate struct {
    dingtalk.BaseCorp
    corpId         string
    agentTag       string
    atType         string
    userId         string // 用户id
    labelList      []int  // 标签列表
    followerUserId string // 负责人用户id
    name           string // 名称
}

func (ecu *extContactUpdate) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        ecu.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不合法", nil))
    }
}

func (ecu *extContactUpdate) SetLabelList(labelList []int) {
    ecu.labelList = make([]int, 0)
    for _, v := range labelList {
        if v > 0 {
            ecu.labelList = append(ecu.labelList, v)
        }
    }

    if len(ecu.labelList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "标签列表不能为空", nil))
    }
}

func (ecu *extContactUpdate) SetFollowerUserId(followerUserId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, followerUserId)
    if match {
        ecu.followerUserId = followerUserId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "负责人用户id不合法", nil))
    }
}

func (ecu *extContactUpdate) SetName(name string) {
    if len(name) > 0 {
        ecu.name = name
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "名称不合法", nil))
    }
}

func (ecu *extContactUpdate) SetTitle(title string) {
    ecu.ExtendData["title"] = strings.TrimSpace(title)
}

func (ecu *extContactUpdate) SetShareDeptIds(shareDeptIds []int) {
    deptList := make([]int, 0)
    for _, v := range shareDeptIds {
        if v > 0 {
            deptList = append(deptList, v)
        }
    }
    if len(deptList) > 0 {
        ecu.ExtendData["share_dept_ids"] = deptList
    }
}

func (ecu *extContactUpdate) SetAddress(address string) {
    ecu.ExtendData["address"] = strings.TrimSpace(address)
}

func (ecu *extContactUpdate) SetRemark(remark string) {
    ecu.ExtendData["remark"] = strings.TrimSpace(remark)
}

func (ecu *extContactUpdate) SetCompanyName(companyName string) {
    ecu.ExtendData["company_name"] = strings.TrimSpace(companyName)
}

func (ecu *extContactUpdate) SetShareUserIds(shareUserIds []string) {
    userList := make([]string, 0)
    for _, v := range shareUserIds {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            userList = append(userList, v)
        }
    }
    if len(userList) > 0 {
        ecu.ExtendData["share_user_ids"] = userList
    }
}

func (ecu *extContactUpdate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ecu.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不能为空", nil))
    }
    if len(ecu.labelList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "标签列表不能为空", nil))
    }
    if len(ecu.followerUserId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "负责人用户id不能为空", nil))
    }
    if len(ecu.name) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "名称不能为空", nil))
    }
    ecu.ExtendData["user_id"] = ecu.userId
    ecu.ExtendData["label_ids"] = ecu.labelList
    ecu.ExtendData["follower_user_id"] = ecu.followerUserId
    ecu.ExtendData["name"] = ecu.name

    ecu.ReqUrl = dingtalk.UrlService + "/topapi/extcontact/update?access_token=" + dingtalk.NewUtil().GetAccessToken(ecu.corpId, ecu.agentTag, ecu.atType)

    reqData := make(map[string]interface{})
    reqData["contact"] = ecu.ExtendData
    reqBody := mpf.JSONMarshal(reqData)
    client, req := ecu.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewExtContactUpdate(corpId, agentTag, atType string) *extContactUpdate {
    ecu := &extContactUpdate{dingtalk.NewCorp(), "", "", "", "", make([]int, 0), "", ""}
    ecu.corpId = corpId
    ecu.agentTag = agentTag
    ecu.atType = atType
    ecu.ReqContentType = project.HTTPContentTypeJSON
    ecu.ReqMethod = fasthttp.MethodPost
    return ecu
}
