package media

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 分块上传事务
type transactionUpload struct {
    dingtalk.BaseCorp
    corpId       string
    agentTag     string
    atType       string
    fileSize     int64 // 文件大小
    chunkNumbers int   // 文件总块数
}

func (tu *transactionUpload) SetFileSize(fileSize int64) {
    if fileSize >= 102400 {
        tu.fileSize = fileSize
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件大小不合法", nil))
    }
}

func (tu *transactionUpload) SetChunkNumbers(chunkNumbers int) {
    if (chunkNumbers > 0) && (chunkNumbers <= 10000) {
        tu.chunkNumbers = chunkNumbers
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件总块数不合法", nil))
    }
}

// 上传事务id,没有设置是开启事务,设置了是提交事务
func (tu *transactionUpload) SetUploadId(uploadId string) {
    if len(uploadId) > 0 {
        tu.ReqData["upload_id"] = uploadId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "上传事务id不合法", nil))
    }
}

func (tu *transactionUpload) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if tu.fileSize <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件大小不能为空", nil))
    }
    if tu.chunkNumbers <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件总块数不能为空", nil))
    }
    tu.ReqData["file_size"] = strconv.FormatInt(tu.fileSize, 10)
    tu.ReqData["chunk_numbers"] = strconv.Itoa(tu.chunkNumbers)
    tu.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(tu.corpId, tu.agentTag, tu.atType)
    tu.ReqUrl = dingtalk.UrlService + "/file/upload/transaction?" + mpf.HttpCreateParams(tu.ReqData, "none", 1)

    return tu.GetRequest()
}

func NewTransactionUpload(corpId, agentTag, atType string) *transactionUpload {
    agentInfo := dingtalk.NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    tu := &transactionUpload{dingtalk.NewCorp(), "", "", "", 0, 0}
    tu.corpId = corpId
    tu.agentTag = agentTag
    tu.atType = atType
    tu.ReqData["agent_id"] = agentInfo["id"]
    return tu
}
