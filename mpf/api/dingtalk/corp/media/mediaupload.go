package media

import (
    "bytes"
    "io"
    "mime/multipart"
    "os"
    "path"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 上传媒体文件
type mediaUpload struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    fileType string
    filePath string // 文件全路径,包括文件名
}

func (mu *mediaUpload) SetFileType(fileType string) {
    _, ok := dingtalk.MediaTypes[fileType]
    if ok {
        mu.fileType = fileType
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件类型不合法", nil))
    }
}

func (mu *mediaUpload) SetFilePath(filePath string) {
    f, err := os.Stat(filePath)
    if err != nil {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不能是目录", nil))
    }
    mu.filePath = filePath
}

func (mu *mediaUpload) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(mu.fileType) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件类型不能为空", nil))
    }
    if len(mu.filePath) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不能为空", nil))
    }
    mu.ReqData["type"] = mu.fileType
    mu.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(mu.corpId, mu.agentTag, mu.atType)
    mu.ReqUrl = dingtalk.UrlService + "/media/upload?" + mpf.HttpCreateParams(mu.ReqData, "none", 1)

    // 新建一个缓冲，用于存放文件内容
    bodyBuffer := &bytes.Buffer{}
    // 创建一个multipart文件写入器,方便按照http规定格式写入内容
    bodyWriter := multipart.NewWriter(bodyBuffer)
    defer bodyWriter.Close()
    // 从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次
    fileWriter, err := bodyWriter.CreateFormFile("media", path.Base(mu.filePath))
    if err != nil {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不合法", nil))
    }

    file, err := os.Open(mu.filePath)
    defer file.Close()
    if err != nil {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不合法", nil))
    }

    _, err = io.Copy(fileWriter, file)
    if err != nil {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不合法", nil))
    }

    mu.ReqContentType = bodyWriter.FormDataContentType()
    client, req := mu.GetRequest()
    req.SetBody(bodyBuffer.Bytes())

    return client, req
}

func NewMediaUpload(corpId, agentTag, atType string) *mediaUpload {
    mu := &mediaUpload{dingtalk.NewCorp(), "", "", "", "", ""}
    mu.corpId = corpId
    mu.agentTag = agentTag
    mu.atType = atType
    mu.ReqMethod = fasthttp.MethodPost
    return mu
}
