/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 16:23
 */
package ota

import (
    "bytes"
    "io"
    "mime/multipart"
    "os"
    "path"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 上传固件包文件
type firmwareFile struct {
    mpiot.BaseBaiDu
    filePath string // 文件全路径,包括文件名
}

func (ff *firmwareFile) SetFilePath(filePath string) {
    f, err := os.Stat(filePath)
    if err != nil {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "文件不能是目录", nil))
    }
    ff.filePath = filePath
}

func (ff *firmwareFile) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ff.filePath) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "文件不能为空", nil))
    }

    // 新建一个缓冲，用于存放文件内容
    bodyBuffer := &bytes.Buffer{}
    // 创建一个multipart文件写入器,方便按照http规定格式写入内容
    bodyWriter := multipart.NewWriter(bodyBuffer)
    defer bodyWriter.Close()
    // 从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次
    fileWriter, err := bodyWriter.CreateFormFile("file", path.Base(ff.filePath))
    if err != nil {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "文件不合法", nil))
    }

    file, err := os.Open(ff.filePath)
    defer file.Close()
    if err != nil {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "文件不合法", nil))
    }

    _, err = io.Copy(fileWriter, file)
    if err != nil {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "文件不合法", nil))
    }

    ff.ReqUrl = ff.GetServiceUrl()
    ff.ReqContentType = bodyWriter.FormDataContentType()

    client, req := ff.GetRequest()
    req.SetBody([]byte(bodyBuffer.Bytes()))

    return client, req
}

func NewFirmwareFile() *firmwareFile {
    ff := &firmwareFile{mpiot.NewBaseBaiDu(), ""}
    ff.ServiceUri = "/v3/iot/management/ota/firmware-file"
    ff.ReqHeader["Content-Type"] = "multipart/form-data"
    ff.ReqMethod = fasthttp.MethodPost
    return ff
}
