/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 13:04
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

// 更新成员
type userUpdate struct {
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

func (uu *userUpdate) SetUserId(userId string) {
    match, _ := regexp.MatchString(`^[0-9a-z]{1,32}$`, userId)
    if match {
        uu.userId = userId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户ID不合法", nil))
    }
}

func (uu *userUpdate) SetName(name string) {
    if len(name) > 0 {
        trueName := []rune(name)
        uu.name = string(trueName[:32])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "名称不合法", nil))
    }
}

func (uu *userUpdate) SetAlias(alias string) {
    if len(alias) > 0 {
        trueAlias := []rune(alias)
        uu.alias = string(trueAlias[:16])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "别名不合法", nil))
    }
}

func (uu *userUpdate) SetMobile(mobile string) {
    match, _ := regexp.MatchString(`^[0-9]{11}$`, mobile)
    if match && (mobile[:1] == "1") {
        uu.mobile = mobile
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "手机号码不合法", nil))
    }
}

func (uu *userUpdate) SetDepartmentList(departmentList []map[string]int) {
    if len(departmentList) > 20 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门列表不能超过20个", nil))
    }

    uu.departments = make([]int, 0)
    uu.orders = make([]int, 0)
    uu.leaders = make([]int, 0)
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
        uu.departments = append(uu.departments, v["depart_id"])
        uu.orders = append(uu.orders, v["order_num"])
        uu.leaders = append(uu.leaders, v["leader_flag"])
    }
}

func (uu *userUpdate) SetPosition(position string) {
    truePosition := []rune(position)
    uu.position = string(truePosition[:64])
}

func (uu *userUpdate) SetGender(gender int) {
    if (gender == 1) || (gender == 2) {
        uu.gender = gender
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "性别不合法", nil))
    }
}

func (uu *userUpdate) SetEmail(email string) {
    if len(email) > 0 {
        uu.email = email
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "邮箱不合法", nil))
    }
}

func (uu *userUpdate) SetTelephone(telephone string) {
    match, _ := regexp.MatchString(`^[0-9\-]+$`, telephone)
    if match {
        uu.telephone = telephone
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "座机不合法", nil))
    }
}

func (uu *userUpdate) SetAvatarMediaId(avatarMediaId string) {
    if len(avatarMediaId) > 0 {
        uu.avatarMediaId = avatarMediaId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "头像不合法", nil))
    }
}

func (uu *userUpdate) SetEnableFlag(enableFlag int) {
    if (enableFlag == 0) || (enableFlag == 1) {
        uu.enableFlag = enableFlag
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "成员标识不合法", nil))
    }
}

func (uu *userUpdate) SetExtArr(extArr map[string]interface{}) {
    if len(extArr) > 0 {
        uu.extArr = extArr
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "扩展属性不合法", nil))
    }
}

func (uu *userUpdate) SetInviteFlag(inviteFlag bool) {
    uu.inviteFlag = inviteFlag
}

func (uu *userUpdate) SetExternalProfile(externalProfile map[string]interface{}) {
    if len(externalProfile) > 0 {
        uu.externalProfile = externalProfile
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "对外属性不合法", nil))
    }
}

func (uu *userUpdate) SetExternalPosition(externalPosition string) {
    if len(externalPosition) > 0 {
        truePosition := []rune(externalPosition)
        uu.externalPosition = string(truePosition[:12])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "对外职务不合法", nil))
    }
}

func (uu *userUpdate) checkData() {
    if len(uu.userId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户ID不能为空", nil))
    }
}

func (uu *userUpdate) SendRequest(getType string) api.ApiResult {
    uu.checkData()

    reqData := make(map[string]interface{})
    reqData["userid"] = uu.userId
    reqData["gender"] = uu.gender
    reqData["enable"] = uu.enableFlag
    reqData["to_invite"] = uu.inviteFlag
    if len(uu.departments) > 0 {
        reqData["department"] = uu.departments
        reqData["order"] = uu.orders
        reqData["is_leader_in_dept"] = uu.leaders
    }
    if len(uu.name) > 0 {
        reqData["name"] = uu.name
    }
    if len(uu.position) > 0 {
        reqData["position"] = uu.position
    }
    if len(uu.alias) > 0 {
        reqData["alias"] = uu.alias
    }
    if len(uu.email) > 0 {
        reqData["email"] = uu.email
    }
    if len(uu.mobile) > 0 {
        reqData["mobile"] = uu.mobile
    }
    if len(uu.telephone) > 0 {
        reqData["telephone"] = uu.telephone
    }
    if len(uu.avatarMediaId) > 0 {
        reqData["avatar_mediaid"] = uu.avatarMediaId
    }
    if len(uu.extArr) > 0 {
        reqData["extattr"] = uu.extArr
    }
    if len(uu.externalProfile) > 0 {
        reqData["external_profile"] = uu.externalProfile
    }
    if len(uu.externalPosition) > 0 {
        reqData["external_position"] = uu.externalPosition
    }
    reqBody := mpf.JsonMarshal(reqData)
    uu.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/user/update?access_token=" + wx.NewUtilWx().GetCorpCache(uu.corpId, uu.agentTag, getType)
    client, req := uu.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := uu.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewUserUpdate(corpId, agentTag string) *userUpdate {
    uu := &userUpdate{wx.NewBaseWxCorp(), "", "", "", "", "", "", make([]int, 0), make([]int, 0), make([]int, 0), "", 0, "", "", "", 0, make(map[string]interface{}), false, make(map[string]interface{}), ""}
    uu.corpId = corpId
    uu.agentTag = agentTag
    uu.gender = 1
    uu.enableFlag = 1
    uu.inviteFlag = true
    uu.ReqContentType = project.HttpContentTypeJson
    uu.ReqMethod = fasthttp.MethodPost
    return uu
}
