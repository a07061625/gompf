/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 16:04
 */
package user

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpim"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type profileSet struct {
    mpim.BaseTencent
    userId    string                   // 用户ID
    userItems []map[string]interface{} // 用户资料列表
}

func (ps *profileSet) SetUserId(userId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, userId)
    if match {
        ps.userId = userId
    } else {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "用户ID不合法", nil))
    }
}

func (ps *profileSet) SetUserItems(userItems []map[string]interface{}) {
    if len(userItems) == 0 {
        ps.userItems = userItems
    } else {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "用户资料列表不合法", nil))
    }
}

func (ps *profileSet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ps.userId) == 0 {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "用户ID不能为空", nil))
    }
    if len(ps.userItems) == 0 {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "用户资料不能为空", nil))
    }
    ps.ExtendData["From_Account"] = ps.userId
    ps.ExtendData["ProfileItem"] = ps.userItems
    reqBody := mpf.JsonMarshal(ps.ExtendData)

    client, req := ps.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewProfileSet() *profileSet {
    ps := &profileSet{mpim.NewBaseTencent(), "", make([]map[string]interface{}, 0)}
    ps.ServiceUri = "/profile/portrait_set"
    ps.ReqContentType = project.HTTPContentTypeJSON
    ps.ReqMethod = fasthttp.MethodPost
    return ps
}
