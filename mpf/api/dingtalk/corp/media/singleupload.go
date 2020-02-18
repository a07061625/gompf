package media

import (
    "bytes"
    "io"
    "mime/multipart"
    "os"
    "path"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 单步文件上传
type singleUpload struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    filePath string // 文件全路径,包括文件名
    fileSize int64  // 文件大小
}

func (su *singleUpload) SetFilePath(filePath string) {
    f, err := os.Stat(filePath)
    if err != nil {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不能是目录", nil))
    }
    su.filePath = filePath
    su.fileSize = f.Size()
}

func (su *singleUpload) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(su.filePath) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不能为空", nil))
    }
    su.ReqData["file_size"] = strconv.FormatInt(su.fileSize, 10)
    su.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(su.corpId, su.agentTag, su.atType)
    su.ReqUrl = dingtalk.UrlService + "/file/upload/single?" + mpf.HTTPCreateParams(su.ReqData, "none", 1)

    // 新建一个缓冲，用于存放文件内容
    bodyBuffer := &bytes.Buffer{}
    // 创建一个multipart文件写入器,方便按照http规定格式写入内容
    bodyWriter := multipart.NewWriter(bodyBuffer)
    defer bodyWriter.Close()
    // 从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次
    fileWriter, err := bodyWriter.CreateFormFile("file", path.Base(su.filePath))
    if err != nil {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不合法", nil))
    }

    file, err := os.Open(su.filePath)
    defer file.Close()
    if err != nil {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不合法", nil))
    }

    _, err = io.Copy(fileWriter, file)
    if err != nil {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不合法", nil))
    }

    su.ReqContentType = bodyWriter.FormDataContentType()
    client, req := su.GetRequest()
    req.SetBody(bodyBuffer.Bytes())

    return client, req
}

func NewSingleUpload(corpId, agentTag, atType string) *singleUpload {
    agentInfo := dingtalk.NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    su := &singleUpload{dingtalk.NewCorp(), "", "", "", "", 0}
    su.corpId = corpId
    su.agentTag = agentTag
    su.atType = atType
    su.ReqData["agent_id"] = agentInfo["id"]
    su.ReqMethod = fasthttp.MethodPost
    return su
}
