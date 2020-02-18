/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 16:19
 */
package pay

import (
    "crypto/tls"
    "encoding/xml"
    "io/ioutil"
    "os"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type commentQueryBatch struct {
    wx.BaseWxAccount
    appId      string
    beginTime  string // 开始时间
    endTime    string // 结束时间
    offset     int    // 位移
    limit      int    // 条数
    outputFile string // 输出文件全名
}

func (cqb *commentQueryBatch) SetTime(beginTime, endTime int) {
    if beginTime <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "开始时间不合法", nil))
    } else if endTime <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "结束时间不合法", nil))
    } else if beginTime > endTime {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "结束时间不能小于开始时间", nil))
    }

    bt := time.Unix(int64(beginTime), 0)
    et := time.Unix(int64(endTime), 0)
    cqb.beginTime = bt.Format("20060102030405")
    cqb.endTime = et.Format("20060102030405")
}

func (cqb *commentQueryBatch) SetOffset(offset int) {
    if offset >= 0 {
        cqb.offset = offset
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "位移不合法", nil))
    }
}

func (cqb *commentQueryBatch) SetLimit(limit int) {
    if (limit > 0) && (limit <= 200) {
        cqb.limit = limit
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "条数不合法", nil))
    }
}

func (cqb *commentQueryBatch) SetOutputFile(outputFile string) {
    f, err := os.Stat(outputFile)
    if err != nil {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "输出文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "输出文件不合法", nil))
    }
    cqb.outputFile = outputFile
}

func (cqb *commentQueryBatch) checkData() {
    if len(cqb.beginTime) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "开始时间不能为空", nil))
    }
    if len(cqb.outputFile) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "输出文件不能为空", nil))
    }
    cqb.ReqData["begin_time"] = cqb.beginTime
    cqb.ReqData["end_time"] = cqb.endTime
    cqb.ReqData["offset"] = strconv.Itoa(cqb.offset)
    cqb.ReqData["limit"] = strconv.Itoa(cqb.limit)
}

func (cqb *commentQueryBatch) SendRequest() api.ApiResult {
    cqb.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(cqb.ReqData, cqb.appId, "md5")
    cqb.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(cqb.ReqData))
    cqb.ReqUrl = "https://api.mch.weixin.qq.com/billcommentsp/batchquerycomment"
    client, req := cqb.GetRequest()
    req.SetBody([]byte(reqBody))

    conf := wx.NewConfig().GetAccount(cqb.appId)
    keyFile, _ := ioutil.TempFile("", "tmpfile")
    defer os.Remove(keyFile.Name())
    if _, err := keyFile.Write([]byte(conf.GetSslKey())); err != nil {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "写入证书key文件失败", nil))
    }

    certFile, _ := ioutil.TempFile("", "tmpfile")
    defer os.Remove(certFile.Name())
    if _, err := certFile.Write([]byte(conf.GetSslCert())); err != nil {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "写入证书cert文件失败", nil))
    }

    certs, err := tls.LoadX509KeyPair(certFile.Name(), keyFile.Name())
    if err != nil {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "加载证书文件失败", nil))
    }
    client.TLSConfig.Certificates = []tls.Certificate{certs}

    resp, result := cqb.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData := make(map[string]string)
    err = xml.Unmarshal(resp.Body, (*mpf.XMLMap)(&respData))
    if err != nil {
        f, err := os.Create(cqb.outputFile)
        defer f.Close()
        if err != nil {
            result.Code = errorcode.WxAccountRequestPost
            result.Msg = err.Error()
        } else {
            f.Write(resp.Body)
            resultData := make(map[string]string)
            resultData["return_code"] = "SUCCESS"
            result.Data = resultData
        }
    } else if respData["return_code"] == "FAIL" {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["return_msg"]
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["err_code_des"]
    }
    return result
}

func NewCommentQueryBatch(appId string) *commentQueryBatch {
    conf := wx.NewConfig().GetAccount(appId)
    cqb := &commentQueryBatch{wx.NewBaseWxAccount(), "", "", "", 0, 0, ""}
    cqb.appId = appId
    cqb.offset = 0
    cqb.limit = 100
    cqb.ReqData["appid"] = conf.GetAppId()
    cqb.ReqData["mch_id"] = conf.GetPayMchId()
    cqb.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    cqb.ReqData["sign_type"] = "HMAC-SHA256"
    cqb.ReqContentType = project.HTTPContentTypeXML
    cqb.ReqMethod = fasthttp.MethodPost
    return cqb
}
