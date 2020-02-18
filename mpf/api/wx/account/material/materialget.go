/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 10:00
 */
package material

import (
    "os"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 下载永久素材
type materialGet struct {
    wx.BaseWxAccount
    appId     string
    mediaId   string // 媒体文件ID
    outputDir string // 输出目录
}

func (md *materialGet) SetMediaId(mediaId string) {
    if len(mediaId) > 0 {
        md.mediaId = mediaId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "媒体文件ID不合法", nil))
    }
}

func (mg *materialGet) SetOutputDir(outputDir string) {
    f, err := os.Stat(outputDir)
    if err != nil {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "输出目录不合法", nil))
    }
    if !f.IsDir() {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "输出目录不合法", nil))
    }
    mg.outputDir = outputDir
}

func (mg *materialGet) checkData() {
    if len(mg.mediaId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "媒体文件ID不能为空", nil))
    }
    if len(mg.outputDir) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "输出目录不能为空", nil))
    }
    mg.ReqData["media_id"] = mg.mediaId
}

func (mg *materialGet) SendRequest() api.ApiResult {
    mg.checkData()

    reqBody := mpf.JsonMarshal(mg.ReqData)
    mg.ReqUrl = "https://api.weixin.qq.com/cgi-bin/material/get_material?access_token=" + wx.NewUtilWx().GetSingleAccessToken(mg.appId)
    client, req := mg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := mg.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, err := mpf.JsonUnmarshalMap(resp.Content)
    if err != nil {
        fileName := mg.outputDir + mg.mediaId
        f, err := os.Create(fileName)
        defer f.Close()
        if err != nil {
            result.Code = errorcode.WxAccountRequestPost
            result.Msg = err.Error()
        } else {
            f.Write(resp.Body)
            resultData := make(map[string]string)
            resultData["media_path"] = fileName
            result.Data = resultData
        }
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewMaterialGet(appId string) *materialGet {
    mg := &materialGet{wx.NewBaseWxAccount(), "", "", ""}
    mg.appId = appId
    mg.ReqContentType = project.HTTPContentTypeJSON
    mg.ReqMethod = fasthttp.MethodPost
    return mg
}
