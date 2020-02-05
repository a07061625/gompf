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

// 上传升级设备文件
type deviceFile struct {
    mpiot.BaseBaiDu
    filePath string // 文件全路径,包括文件名
}

func (df *deviceFile) SetFilePath(filePath string) {
    f, err := os.Stat(filePath)
    if err != nil {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "文件不能是目录", nil))
    }
    df.filePath = filePath
}

func (df *deviceFile) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(df.filePath) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "文件不能为空", nil))
    }

    // 新建一个缓冲，用于存放文件内容
    bodyBuffer := &bytes.Buffer{}
    // 创建一个multipart文件写入器,方便按照http规定格式写入内容
    bodyWriter := multipart.NewWriter(bodyBuffer)
    defer bodyWriter.Close()
    // 从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次
    fileWriter, err := bodyWriter.CreateFormFile("file", path.Base(df.filePath))
    if err != nil {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "文件不合法", nil))
    }

    file, err := os.Open(df.filePath)
    defer file.Close()
    if err != nil {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "文件不合法", nil))
    }

    _, err = io.Copy(fileWriter, file)
    if err != nil {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "文件不合法", nil))
    }

    df.ReqUrl = df.GetServiceUrl()
    df.ReqContentType = bodyWriter.FormDataContentType()

    client, req := df.GetRequest()
    req.SetBody([]byte(bodyBuffer.Bytes()))

    return client, req
}

func NewDeviceFile() *deviceFile {
    df := &deviceFile{mpiot.NewBaseBaiDu(), ""}
    df.ServiceUri = "/v3/iot/management/ota/device-file"
    df.ReqHeader["Content-Type"] = "multipart/form-data"
    df.ReqMethod = fasthttp.MethodPost
    return df
}
