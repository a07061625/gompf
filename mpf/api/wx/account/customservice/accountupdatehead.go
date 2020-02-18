/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 21:57
 */
package customservice

import (
    "bytes"
    "io"
    "mime/multipart"
    "os"
    "path"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type accountUpdateHead struct {
    wx.BaseWxAccount
    appId     string
    kfAccount string // 客服帐号 格式为: 帐号前缀@公众号微信号
    filePath  string // 文件全路径,包括文件名
}

func (auh *accountUpdateHead) SetKfAccount(kfAccount string) {
    if (len(kfAccount) > 0) && (len(kfAccount) <= 30) {
        auh.kfAccount = kfAccount
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服帐号不合法", nil))
    }
}

func (auh *accountUpdateHead) SetFilePath(filePath string) {
    f, err := os.Stat(filePath)
    if err != nil {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文件不能是目录", nil))
    }
    auh.filePath = filePath
}

func (auh *accountUpdateHead) checkData() (string, []byte, error) {
    if len(auh.kfAccount) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服帐号不能为空", nil))
    }
    if len(auh.filePath) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文件不能为空", nil))
    }
    auh.ReqData["kf_account"] = auh.kfAccount

    // 新建一个缓冲，用于存放文件内容
    bodyBuffer := &bytes.Buffer{}
    // 创建一个multipart文件写入器,方便按照http规定格式写入内容
    bodyWriter := multipart.NewWriter(bodyBuffer)
    defer bodyWriter.Close()
    // 从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次
    fileWriter, err := bodyWriter.CreateFormFile("media", path.Base(auh.filePath))
    if err != nil {
        return "", nil, err
    }

    file, err := os.Open(auh.filePath)
    defer file.Close()
    if err != nil {
        return "", nil, err
    }

    _, err = io.Copy(fileWriter, file)
    if err != nil {
        return "", nil, err
    }

    contentType := bodyWriter.FormDataContentType()
    content := bodyBuffer.Bytes()
    return contentType, content, nil
}

func (auh *accountUpdateHead) SendRequest() api.ApiResult {
    contentType, content, err := auh.checkData()
    if err != nil {
        result := api.NewApiResult()
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = "上传文件失败"
        return result
    }

    auh.ReqUrl = "https://api.weixin.qq.com/customservice/kfaccount/uploadheadimg?access_token=" + wx.NewUtilWx().GetSingleAccessToken(auh.appId)
    client, req := auh.GetRequest()
    req.Header.SetContentType(contentType)
    req.SetBody(content)

    resp, result := auh.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewAccountUpdateHead(appId string) *accountUpdateHead {
    auh := &accountUpdateHead{wx.NewBaseWxAccount(), "", "", ""}
    auh.appId = appId
    auh.ReqContentType = project.HTTPContentTypeJSON
    auh.ReqMethod = fasthttp.MethodPost
    return auh
}
