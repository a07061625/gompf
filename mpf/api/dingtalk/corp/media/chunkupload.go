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

// 上传文件块
type chunkUpload struct {
    dingtalk.BaseCorp
    corpId        string
    agentTag      string
    atType        string
    filePath      string // 文件全路径,包括文件名
    uploadId      string // 上传事务id
    chunkSequence int    // 文件块号
}

func (cu *chunkUpload) SetFilePath(filePath string) {
    f, err := os.Stat(filePath)
    if err != nil {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不能是目录", nil))
    }
    cu.filePath = filePath
}

func (cu *chunkUpload) SetUploadId(uploadId string) {
    if len(uploadId) > 0 {
        cu.uploadId = uploadId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "上传事务id不合法", nil))
    }
}

func (cu *chunkUpload) SetChunkSequence(chunkSequence int) {
    if chunkSequence > 0 {
        cu.chunkSequence = chunkSequence
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件块号不合法", nil))
    }
}

func (cu *chunkUpload) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(cu.filePath) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不能为空", nil))
    }
    if len(cu.uploadId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "上传事务id不能为空", nil))
    }
    if cu.chunkSequence == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件块号不能为空", nil))
    }
    cu.ReqData["upload_id"] = cu.uploadId
    cu.ReqData["chunk_sequence"] = strconv.Itoa(cu.chunkSequence)
    cu.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(cu.corpId, cu.agentTag, cu.atType)
    cu.ReqUrl = dingtalk.UrlService + "/file/upload/chunk?" + mpf.HTTPCreateParams(cu.ReqData, "none", 1)

    // 新建一个缓冲，用于存放文件内容
    bodyBuffer := &bytes.Buffer{}
    // 创建一个multipart文件写入器,方便按照http规定格式写入内容
    bodyWriter := multipart.NewWriter(bodyBuffer)
    defer bodyWriter.Close()
    // 从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次
    fileWriter, err := bodyWriter.CreateFormFile("file", path.Base(cu.filePath))
    if err != nil {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不合法", nil))
    }

    file, err := os.Open(cu.filePath)
    defer file.Close()
    if err != nil {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不合法", nil))
    }

    _, err = io.Copy(fileWriter, file)
    if err != nil {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件不合法", nil))
    }

    cu.ReqContentType = bodyWriter.FormDataContentType()
    client, req := cu.GetRequest()
    req.SetBody(bodyBuffer.Bytes())

    return client, req
}

func NewChunkUpload(corpId, agentTag, atType string) *chunkUpload {
    agentInfo := dingtalk.NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    cu := &chunkUpload{dingtalk.NewCorp(), "", "", "", "", "", 0}
    cu.corpId = corpId
    cu.agentTag = agentTag
    cu.atType = atType
    cu.ReqData["agent_id"] = agentInfo["id"]
    cu.ReqMethod = fasthttp.MethodPost
    return cu
}
