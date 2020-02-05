/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 10:08
 */
package material

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

// 新增其他类型永久素材
type materialAdd struct {
    wx.BaseWxAccount
    appId        string
    materialType string                 // 素材类型
    fileInfo     map[string]interface{} // 文件信息
    filePath     string                 // 文件全路径,包括文件名
}

func (ma *materialAdd) SetMaterialType(materialType string) {
    _, ok := account.MaterialTypes[materialType]
    if ok {
        ma.materialType = materialType
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "素材类型不合法", nil))
    }
}

func (ma *materialAdd) SetFileInfo(fileInfo map[string]interface{}) {
    if len(fileInfo) == 0 {
        ma.fileInfo = fileInfo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文件信息不合法", nil))
    }
}

func (ma *materialAdd) SetFilePath(filePath string) {
    f, err := os.Stat(filePath)
    if err != nil {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文件不能是目录", nil))
    }
    ma.filePath = filePath
}

func (ma *materialAdd) checkData() (string, []byte, error) {
    if len(ma.materialType) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "素材类型不能为空", nil))
    }
    if len(ma.fileInfo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文件信息不能为空", nil))
    }
    if len(ma.filePath) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文件不能为空", nil))
    }
    fileInfoStr := mpf.JsonMarshal(ma.fileInfo)

    // 新建一个缓冲，用于存放文件内容
    bodyBuffer := &bytes.Buffer{}
    // 创建一个multipart文件写入器,方便按照http规定格式写入内容
    bodyWriter := multipart.NewWriter(bodyBuffer)
    defer bodyWriter.Close()

    extWriter, _ := bodyWriter.CreateFormField("description")
    extWriter.Write([]byte(fileInfoStr))

    // 从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次
    fileWriter, err := bodyWriter.CreateFormFile("media", path.Base(ma.filePath))
    if err != nil {
        return "", nil, err
    }

    file, err := os.Open(ma.filePath)
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

func (ma *materialAdd) SendRequest() api.ApiResult {
    contentType, content, err := ma.checkData()
    if err != nil {
        result := api.NewApiResult()
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = "上传文件失败"
        return result
    }

    ma.ReqData["type"] = ma.materialType
    ma.ReqData["access_token"] = wx.NewUtilWx().GetSingleAccessToken(ma.appId)
    ma.ReqUrl = "https://api.weixin.qq.com/cgi-bin/material/add_material?" + mpf.HttpCreateParams(ma.ReqData, "none", 1)
    client, req := ma.GetRequest()
    req.Header.SetContentType(contentType)
    req.SetBody(content)

    resp, result := ma.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["media_id"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewMaterialAdd(appId string) *materialAdd {
    ma := &materialAdd{wx.NewBaseWxAccount(), "", "", make(map[string]interface{}), ""}
    ma.appId = appId
    ma.ReqMethod = fasthttp.MethodPost
    return ma
}
