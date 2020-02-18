/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 21:57
 */
package auth

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 认证
type authenticatePassword struct {
    mpiot.BaseBaiDu
    userName string // thing的用户名
    password string // 身份密钥
}

func (ap *authenticatePassword) SetUserName(userName string) {
    if len(userName) > 0 {
        ap.userName = userName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "用户名不合法", nil))
    }
}

func (ap *authenticatePassword) SetPassword(password string) {
    if len(password) > 0 {
        ap.password = password
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "身份密钥不合法", nil))
    }
}

func (ap *authenticatePassword) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ap.userName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "用户名不能为空", nil))
    }
    if len(ap.password) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "身份密钥不能为空", nil))
    }
    ap.ExtendData["username"] = ap.userName
    ap.ExtendData["password"] = ap.password

    ap.ReqUrl = ap.GetServiceUrl()

    reqBody := mpf.JsonMarshal(ap.ExtendData)
    client, req := ap.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewAuthenticatePassword() *authenticatePassword {
    ap := &authenticatePassword{mpiot.NewBaseBaiDu(), "", ""}
    ap.ServiceUri = "/v1/auth/authenticate/password"
    ap.ReqContentType = project.HTTPContentTypeJSON
    ap.ReqMethod = fasthttp.MethodPost
    return ap
}
