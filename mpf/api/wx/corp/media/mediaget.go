/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 14:04
 */
package media

import (
    "os"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 获取临时素材
type mediaGet struct {
    wx.BaseWxCorp
    corpId    string
    agentTag  string
    outputDir string // 输出目录
    mediaId   string // 媒体文件ID
}

func (mg *mediaGet) SetOutputDir(outputDir string) {
    f, err := os.Stat(outputDir)
    if err != nil {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "输出目录不合法", nil))
    }
    if !f.IsDir() {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "输出目录不合法", nil))
    }
    mg.outputDir = outputDir
}

func (mg *mediaGet) SetMediaId(mediaId string) {
    if len(mediaId) > 0 {
        mg.mediaId = mediaId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "媒体文件ID不合法", nil))
    }
}

func (mg *mediaGet) checkData() {
    if len(mg.mediaId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "媒体文件ID不能为空", nil))
    }
    if len(mg.outputDir) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "输出目录不能为空", nil))
    }
    mg.ReqData["media_id"] = mg.mediaId
}

func (mg *mediaGet) SendRequest(getType string) api.ApiResult {
    mg.checkData()

    mg.ReqData["access_token"] = wx.NewUtilWx().GetCorpCache(mg.corpId, mg.agentTag, getType)
    mg.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/media/get?" + mpf.HTTPCreateParams(mg.ReqData, "none", 1)
    client, req := mg.GetRequest()

    resp, result := mg.SendInner(client, req, errorcode.WxCorpRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, err := mpf.JSONUnmarshalMap(resp.Content)
    if err != nil {
        fileName := mg.outputDir + mg.mediaId
        f, err := os.Create(fileName)
        defer f.Close()
        if err != nil {
            result.Code = errorcode.WxCorpRequestGet
            result.Msg = err.Error()
        } else {
            f.Write(resp.Body)
            resultData := make(map[string]string)
            resultData["media_path"] = fileName
            result.Data = resultData
        }
    } else {
        result.Code = errorcode.WxCorpRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewMediaGet(corpId, agentTag string) *mediaGet {
    mg := &mediaGet{wx.NewBaseWxCorp(), "", "", "", ""}
    mg.corpId = corpId
    mg.agentTag = agentTag
    return mg
}
