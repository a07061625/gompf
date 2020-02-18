/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 11:04
 */
package media

import (
    "bytes"
    "io"
    "mime/multipart"
    "os"
    "path"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/api/wx/account"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type mediaUpload struct {
    wx.BaseWxAccount
    appId    string
    fileType string // 媒体文件类型
    filePath string // 文件全路径,包括文件名
}

func (mu *mediaUpload) SetFileType(fileType string) {
    _, ok := account.MaterialTypes[fileType]
    if ok {
        mu.fileType = fileType
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "媒体文件类型不合法", nil))
    }
}

func (mu *mediaUpload) SetFilePath(filePath string) {
    f, err := os.Stat(filePath)
    if err != nil {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文件不能是目录", nil))
    }
    mu.filePath = filePath
}

func (mu *mediaUpload) checkData() (string, []byte, error) {
    if len(mu.fileType) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "媒体文件ID不能为空", nil))
    }
    if len(mu.filePath) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文件不能为空", nil))
    }
    mu.ReqData["type"] = mu.fileType

    // 新建一个缓冲，用于存放文件内容
    bodyBuffer := &bytes.Buffer{}
    // 创建一个multipart文件写入器,方便按照http规定格式写入内容
    bodyWriter := multipart.NewWriter(bodyBuffer)
    defer bodyWriter.Close()
    // 从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次
    fileWriter, err := bodyWriter.CreateFormFile("media", path.Base(mu.filePath))
    if err != nil {
        return "", nil, err
    }

    file, err := os.Open(mu.filePath)
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

func (mu *mediaUpload) SendRequest() api.APIResult {
    contentType, content, err := mu.checkData()
    if err != nil {
        result := api.NewAPIResult()
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = "上传文件失败"
        return result
    }

    mu.ReqData["access_token"] = wx.NewUtilWx().GetSingleAccessToken(mu.appId)
    mu.ReqURI = "https://api.weixin.qq.com/cgi-bin/media/upload?" + mpf.HTTPCreateParams(mu.ReqData, "none", 1)
    client, req := mu.GetRequest()
    req.Header.SetContentType(contentType)
    req.SetBody(content)

    resp, result := mu.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["media_id"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewMediaUpload(appId string) *mediaUpload {
    mu := &mediaUpload{wx.NewBaseWxAccount(), "", "", ""}
    mu.appId = appId
    mu.ReqMethod = fasthttp.MethodPost
    return mu
}
