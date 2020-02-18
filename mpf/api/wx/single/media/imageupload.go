/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 14:25
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
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type imageUpload struct {
    wx.BaseWxAccount
    appId    string
    filePath string // 文件全路径,包括文件名
}

func (iu *imageUpload) SetFilePath(filePath string) {
    f, err := os.Stat(filePath)
    if err != nil {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文件不能是目录", nil))
    }
    iu.filePath = filePath
}

func (iu *imageUpload) checkData() (string, []byte, error) {
    if len(iu.filePath) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文件不能为空", nil))
    }

    // 新建一个缓冲，用于存放文件内容
    bodyBuffer := &bytes.Buffer{}
    // 创建一个multipart文件写入器,方便按照http规定格式写入内容
    bodyWriter := multipart.NewWriter(bodyBuffer)
    defer bodyWriter.Close()
    // 从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次
    fileWriter, err := bodyWriter.CreateFormFile("media", path.Base(iu.filePath))
    if err != nil {
        return "", nil, err
    }

    file, err := os.Open(iu.filePath)
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

func (iu *imageUpload) SendRequest() api.APIResult {
    contentType, content, err := iu.checkData()
    if err != nil {
        result := api.NewAPIResult()
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = "上传文件失败"
        return result
    }

    iu.ReqURI = "https://api.weixin.qq.com/cgi-bin/media/uploadimg?access_token=" + wx.NewUtilWx().GetSingleAccessToken(iu.appId)
    client, req := iu.GetRequest()
    req.Header.SetContentType(contentType)
    req.SetBody(content)

    resp, result := iu.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["url"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewImageUpload(appId string) *imageUpload {
    iu := &imageUpload{wx.NewBaseWxAccount(), "", ""}
    iu.appId = appId
    iu.ReqMethod = fasthttp.MethodPost
    return iu
}
