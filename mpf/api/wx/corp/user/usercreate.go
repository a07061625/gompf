/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 13:04
 */
package user

import (
    "regexp"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建成员
type userCreate struct {
    wx.BaseWxCorp
    corpId           string
    agentTag         string
    userId           string                 // 用户ID
    name             string                 // 名称
    alias            string                 // 别名
    mobile           string                 // 手机号码
    departments      []int                  // 部门列表
    orders           []int                  // 排序列表
    leaders          []int                  // 上级列表
    position         string                 // 职务信息
    gender           int                    // 性别 1:男性 2:女性
    email            string                 // 邮箱
    telephone        string                 // 座机
    avatarMediaId    string                 // 头像
    enableFlag       int                    // 成员标识 1:启用 0:禁用
    extArr           map[string]interface{} // 扩展属性
    inviteFlag       bool                   // 邀请标识,默认值为true true:邀请使用企业微信 false:不邀请
    externalProfile  map[string]interface{} // 对外属性
    externalPosition string                 // 对外职务
}

func (uc *userCreate) SetUserId(userId string) {
    match, _ := regexp.MatchString(`^[0-9a-z]{1,32}$`, userId)
    if match {
        uc.userId = userId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户ID不合法", nil))
    }
}

func (uc *userCreate) SetName(name string) {
    if len(name) > 0 {
        trueName := []rune(name)
        uc.name = string(trueName[:32])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "名称不合法", nil))
    }
}

func (uc *userCreate) SetAlias(alias string) {
    if len(alias) > 0 {
        trueAlias := []rune(alias)
        uc.alias = string(trueAlias[:16])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "别名不合法", nil))
    }
}

func (uc *userCreate) SetMobile(mobile string) {
    match, _ := regexp.MatchString(`^[0-9]{11}$`, mobile)
    if match && (mobile[:1] == "1") {
        uc.mobile = mobile
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "手机号码不合法", nil))
    }
}

func (uc *userCreate) SetDepartmentList(departmentList []map[string]int) {
    if len(departmentList) > 20 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门列表不能超过20个", nil))
    }

    uc.departments = make([]int, 0)
    uc.orders = make([]int, 0)
    uc.leaders = make([]int, 0)
    for _, v := range departmentList {
        if v["depart_id"] <= 0 {
            continue
        }
        if v["order_num"] < 0 {
            continue
        }
        if (v["leader_flag"] < 0) || (v["leader_flag"] > 1) {
            continue
        }
        uc.departments = append(uc.departments, v["depart_id"])
        uc.orders = append(uc.orders, v["order_num"])
        uc.leaders = append(uc.leaders, v["leader_flag"])
    }
}

func (uc *userCreate) SetPosition(position string) {
    truePosition := []rune(position)
    uc.position = string(truePosition[:64])
}

func (uc *userCreate) SetGender(gender int) {
    if (gender == 1) || (gender == 2) {
        uc.gender = gender
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "性别不合法", nil))
    }
}

func (uc *userCreate) SetEmail(email string) {
    if len(email) > 0 {
        uc.email = email
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "邮箱不合法", nil))
    }
}

func (uc *userCreate) SetTelephone(telephone string) {
    match, _ := regexp.MatchString(`^[0-9\-]+$`, telephone)
    if match {
        uc.telephone = telephone
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "座机不合法", nil))
    }
}

func (uc *userCreate) SetAvatarMediaId(avatarMediaId string) {
    if len(avatarMediaId) > 0 {
        uc.avatarMediaId = avatarMediaId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "头像不合法", nil))
    }
}

func (uc *userCreate) SetEnableFlag(enableFlag int) {
    if (enableFlag == 0) || (enableFlag == 1) {
        uc.enableFlag = enableFlag
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "成员标识不合法", nil))
    }
}

func (uc *userCreate) SetExtArr(extArr map[string]interface{}) {
    if len(extArr) > 0 {
        uc.extArr = extArr
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "扩展属性不合法", nil))
    }
}

func (uc *userCreate) SetInviteFlag(inviteFlag bool) {
    uc.inviteFlag = inviteFlag
}

func (uc *userCreate) SetExternalProfile(externalProfile map[string]interface{}) {
    if len(externalProfile) > 0 {
        uc.externalProfile = externalProfile
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "对外属性不合法", nil))
    }
}

func (uc *userCreate) SetExternalPosition(externalPosition string) {
    if len(externalPosition) > 0 {
        truePosition := []rune(externalPosition)
        uc.externalPosition = string(truePosition[:12])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "对外职务不合法", nil))
    }
}

func (uc *userCreate) checkData() {
    if len(uc.departments) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门列表不能为空", nil))
    }
    if len(uc.name) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "名称不能为空", nil))
    }
    if (len(uc.mobile) == 0) && (len(uc.email) == 0) {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "手机号码和邮箱不能同时为空", nil))
    }
}

func (uc *userCreate) SendRequest(getType string) api.ApiResult {
    uc.checkData()

    reqData := make(map[string]interface{})
    reqData["department"] = uc.departments
    reqData["order"] = uc.orders
    reqData["is_leader_in_dept"] = uc.leaders
    reqData["name"] = uc.name
    reqData["userid"] = uc.userId
    reqData["gender"] = uc.gender
    reqData["enable"] = uc.enableFlag
    reqData["to_invite"] = uc.inviteFlag
    reqData["position"] = uc.position
    if len(uc.alias) > 0 {
        reqData["alias"] = uc.alias
    }
    if len(uc.email) > 0 {
        reqData["email"] = uc.email
    }
    if len(uc.mobile) > 0 {
        reqData["mobile"] = uc.mobile
    }
    if len(uc.telephone) > 0 {
        reqData["telephone"] = uc.telephone
    }
    if len(uc.avatarMediaId) > 0 {
        reqData["avatar_mediaid"] = uc.avatarMediaId
    }
    if len(uc.extArr) > 0 {
        reqData["extattr"] = uc.extArr
    }
    if len(uc.externalProfile) > 0 {
        reqData["external_profile"] = uc.externalProfile
    }
    if len(uc.externalPosition) > 0 {
        reqData["external_position"] = uc.externalPosition
    }
    reqBody := mpf.JsonMarshal(reqData)
    uc.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/user/create?access_token=" + wx.NewUtilWx().GetCorpCache(uc.corpId, uc.agentTag, getType)
    client, req := uc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := uc.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewUserCreate(corpId, agentTag string) *userCreate {
    uc := &userCreate{wx.NewBaseWxCorp(), "", "", "", "", "", "", make([]int, 0), make([]int, 0), make([]int, 0), "", 0, "", "", "", 0, make(map[string]interface{}), false, make(map[string]interface{}), ""}
    uc.corpId = corpId
    uc.agentTag = agentTag
    uc.userId = mpf.ToolCreateNonceStr(8, "numlower") + strconv.Itoa(time.Now().Second())
    uc.gender = 1
    uc.enableFlag = 1
    uc.inviteFlag = true
    uc.ReqContentType = project.HttpContentTypeJson
    uc.ReqMethod = fasthttp.MethodPost
    return uc
}
