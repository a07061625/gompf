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

// 获取高清语音素材
type voiceGet struct {
    wx.BaseWxCorp
    corpId    string
    agentTag  string
    outputDir string // 输出目录
    mediaId   string // 媒体文件ID
}

func (vg *voiceGet) SetOutputDir(outputDir string) {
    f, err := os.Stat(outputDir)
    if err != nil {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "输出目录不合法", nil))
    }
    if !f.IsDir() {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "输出目录不合法", nil))
    }
    vg.outputDir = outputDir
}

func (vg *voiceGet) SetMediaId(mediaId string) {
    if len(mediaId) > 0 {
        vg.mediaId = mediaId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "媒体文件ID不合法", nil))
    }
}

func (vg *voiceGet) checkData() {
    if len(vg.mediaId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "媒体文件ID不能为空", nil))
    }
    if len(vg.outputDir) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "输出目录不能为空", nil))
    }
    vg.ReqData["media_id"] = vg.mediaId
}

func (vg *voiceGet) SendRequest(getType string) api.APIResult {
    vg.checkData()

    vg.ReqData["access_token"] = wx.NewUtilWx().GetCorpCache(vg.corpId, vg.agentTag, getType)
    vg.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/media/get/jssdkt?" + mpf.HTTPCreateParams(vg.ReqData, "none", 1)
    client, req := vg.GetRequest()

    resp, result := vg.SendInner(client, req, errorcode.WxCorpRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, err := mpf.JSONUnmarshalMap(resp.Content)
    if err != nil {
        fileName := vg.outputDir + vg.mediaId
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

func NewVoiceGet(corpId, agentTag string) *voiceGet {
    mg := &voiceGet{wx.NewBaseWxCorp(), "", "", "", ""}
    mg.corpId = corpId
    mg.agentTag = agentTag
    return mg
}
