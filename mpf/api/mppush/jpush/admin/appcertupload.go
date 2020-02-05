/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 11:17
 */
package admin

import (
    "bytes"
    "io"
    "mime/multipart"
    "os"
    "path"
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// app证书上传
type appCertUpload struct {
    mppush.BaseJPush
    appKey          string // 应用标识
    devCertFile     string // 测试环境证书文件
    devCertPassword string // 测试环境证书密码
    proCertFile     string // 生产环境证书文件
    proCertPassword string // 生产环境证书密码
}

func (acu *appCertUpload) SetAppKey(appKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appKey)
    if match {
        acu.appKey = appKey
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "应用标识不合法", nil))
    }
}

func (acu *appCertUpload) SetDevCert(file, password string) {
    f, err := os.Stat(file)
    if err != nil {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "测试环境证书不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "测试环境证书不能是目录", nil))
    }
    if len(password) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "测试环境证书密码不合法", nil))
    }
    acu.devCertFile = file
    acu.devCertPassword = password
}

func (acu *appCertUpload) SetProCert(file, password string) {
    f, err := os.Stat(file)
    if err != nil {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "生产环境证书不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "生产环境证书不能是目录", nil))
    }
    if len(password) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "生产环境证书密码不合法", nil))
    }
    acu.proCertFile = file
    acu.proCertPassword = password
}

func (acu *appCertUpload) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(acu.appKey) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "应用标识不能为空", nil))
    }
    if (len(acu.devCertFile) == 0) && (len(acu.proCertFile) == 0) {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "生产环境和测试环境证书不能同时为空", nil))
    }
    acu.ServiceUri = "/v1/app/" + acu.appKey + "/certificate"

    // 新建一个缓冲，用于存放文件内容
    bodyBuffer := &bytes.Buffer{}
    // 创建一个multipart文件写入器,方便按照http规定格式写入内容
    bodyWriter := multipart.NewWriter(bodyBuffer)
    defer bodyWriter.Close()
    if len(acu.devCertFile) > 0 {
        devFileWriter, err := bodyWriter.CreateFormFile("devCertificateFile", path.Base(acu.devCertFile))
        if err != nil {
            panic(mperr.NewPushJPush(errorcode.PushJPushParam, "测试环境证书不合法", nil))
        }
        devFile, err := os.Open(acu.devCertFile)
        defer devFile.Close()
        if err != nil {
            panic(mperr.NewPushJPush(errorcode.PushJPushParam, "测试环境证书不合法", nil))
        }

        _, err = io.Copy(devFileWriter, devFile)
        if err != nil {
            panic(mperr.NewPushJPush(errorcode.PushJPushParam, "测试环境证书不合法", nil))
        }

        devPwdWriter, _ := bodyWriter.CreateFormField("devCertificatePassword")
        devPwdWriter.Write([]byte(acu.devCertPassword))
    }
    if len(acu.proCertFile) > 0 {
        proFileWriter, err := bodyWriter.CreateFormFile("proCertificateFile", path.Base(acu.proCertFile))
        if err != nil {
            panic(mperr.NewPushJPush(errorcode.PushJPushParam, "生产环境证书不合法", nil))
        }
        proFile, err := os.Open(acu.proCertFile)
        defer proFile.Close()
        if err != nil {
            panic(mperr.NewPushJPush(errorcode.PushJPushParam, "生产环境证书不合法", nil))
        }

        _, err = io.Copy(proFileWriter, proFile)
        if err != nil {
            panic(mperr.NewPushJPush(errorcode.PushJPushParam, "生产环境证书不合法", nil))
        }

        proPwdWriter, _ := bodyWriter.CreateFormField("proCertificatePassword")
        proPwdWriter.Write([]byte(acu.proCertPassword))
    }

    acu.ReqUrl = acu.GetServiceUrl()
    acu.ReqContentType = bodyWriter.FormDataContentType()
    client, req := acu.GetRequest()
    req.SetBody(bodyBuffer.Bytes())

    return client, req
}

func NewAppCertUpload(key string) *appCertUpload {
    acu := &appCertUpload{mppush.NewBaseJPush(mppush.JPushServiceDomainAdmin, key, "dev"), "", "", "", "", ""}
    acu.ReqHeader["Content-Type"] = "multipart/form-data"
    acu.ReqMethod = fasthttp.MethodPost
    return acu
}
