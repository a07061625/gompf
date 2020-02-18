/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 9:09
 */
package cloud

import (
    "encoding/base64"
    "io/ioutil"
    "os"
    "regexp"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 上传云函数
type functionUpdate struct {
    wx.BaseWxOpen
    appId             string // 应用ID
    handler           string // 函数处理方法名
    functionName      string // 函数名称
    zipFile           string // 函数代码zip文件
    envId             string // 命名空间
    installDependency string // 自动安装依赖标识
}

func (fu *functionUpdate) SetFunctionName(functionName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, functionName)
    if match {
        fu.functionName = functionName
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "函数名称不合法", nil))
    }
}

func (fu *functionUpdate) SetZipFile(zipFile string) {
    f, err := os.Open(zipFile)
    defer f.Close()
    if err != nil {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "函数代码zip文件不合法", nil))
    }
    fc, err := ioutil.ReadAll(f)
    if err != nil {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "函数代码zip文件不合法", nil))
    }
    fu.zipFile = base64.StdEncoding.EncodeToString(fc)
}

func (fu *functionUpdate) SetEnvId(envId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, envId)
    if match {
        fu.envId = envId
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "命名空间不合法", nil))
    }
}

func (fu *functionUpdate) SetInstallDependency(installDependency string) {
    if (installDependency == "TRUE") || (installDependency == "FALSE") {
        fu.ReqData["InstallDependency"] = installDependency
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "自动安装依赖标识不合法", nil))
    }
}

func (fu *functionUpdate) checkData() {
    if len(fu.functionName) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "函数名称不能为空", nil))
    }
    if len(fu.zipFile) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "函数代码zip文件不能为空", nil))
    }
    if len(fu.envId) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "命名空间不能为空", nil))
    }
    fu.ReqData["FunctionName"] = fu.functionName
    fu.ReqData["ZipFile"] = fu.zipFile
    fu.ReqData["EnvId"] = fu.envId
}

func (fu *functionUpdate) SendRequest() api.ApiResult {
    fu.checkData()

    fu.ReqData["CodeSecret"] = wx.NewUtilWx().GetOpenAuthorizeCodeSecret(fu.appId)
    reqBody := mpf.JSONMarshal(fu.ReqData)
    hashedPayload := mpf.HashSha256(reqBody, "")
    signatureGet := NewUploadSignatureGet(fu.appId)
    signatureGet.SetHashedPayload(hashedPayload)
    signRes := signatureGet.SendRequest()
    if signRes.Code > 0 {
        panic(mperr.NewWxOpenMini(signRes.Code, signRes.Msg, nil))
    }

    headerMap := signRes.Data.(map[string]interface{})
    headers := strings.Split(headerMap["headers"].(string), "\r\n")
    for _, v := range headers {
        headerInfo := strings.Split(v, ":")
        fu.ReqHeader[strings.TrimSpace(headerInfo[0])] = strings.TrimSpace(headerInfo[1])
    }
    fu.ReqUrl = "https://scf.tencentcloudapi.com"
    client, req := fu.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := fu.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewFunctionUpdate(appId string) *functionUpdate {
    fu := &functionUpdate{wx.NewBaseWxOpen(), "", "", "", "", "", ""}
    fu.appId = appId
    fu.ReqData["Handler"] = "index.main"
    fu.ReqData["InstallDependency"] = "FALSE"
    fu.ReqContentType = project.HTTPContentTypeJSON
    fu.ReqMethod = fasthttp.MethodPost
    return fu
}
