/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 15:52
 */
package user

import (
    "regexp"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpim"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type accountImport struct {
    mpim.BaseTencent
    userId   string // 用户ID
    userType int    // 用户类型 0:普通用户 1:机器人用户
    nickname string // 昵称
    faceUrl  string // 头像链接
}

func (ai *accountImport) SetUserId(userId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, userId)
    if match {
        ai.userId = userId
    } else {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "用户ID不合法", nil))
    }
}

func (ai *accountImport) SetUserType(userType int) {
    if (userType == 0) || (userType == 1) {
        ai.userType = userType
    } else {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "用户类型不合法", nil))
    }
}

func (ai *accountImport) SetNickname(nickname string) {
    ai.nickname = strings.TrimSpace(nickname)
}

func (ai *accountImport) SetFaceUrl(faceUrl string) {
    match, _ := regexp.MatchString(project.RegexURLHTTP, faceUrl)
    if match {
        ai.faceUrl = faceUrl
    } else {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "用户头像不合法", nil))
    }
}

func (ai *accountImport) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ai.userId) == 0 {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "用户ID不能为空", nil))
    }
    if len(ai.faceUrl) == 0 {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "用户头像不能为空", nil))
    }
    ai.ExtendData["Identifier"] = ai.userId
    ai.ExtendData["Nick"] = ai.nickname
    ai.ExtendData["FaceUrl"] = ai.faceUrl
    ai.ExtendData["Type"] = ai.userType
    reqBody := mpf.JSONMarshal(ai.ExtendData)

    client, req := ai.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewAccountImport() *accountImport {
    ai := &accountImport{mpim.NewBaseTencent(), "", 0, "", ""}
    ai.userType = 0
    ai.ServiceUri = "/im_open_login_svc/account_import"
    ai.ReqContentType = project.HTTPContentTypeJSON
    ai.ReqMethod = fasthttp.MethodPost
    return ai
}
