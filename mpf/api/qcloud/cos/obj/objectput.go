/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 22:03
 */
package obj

import (
    "bytes"
    "encoding/base64"
    "io"
    "mime/multipart"
    "os"
    "path"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 上传本地对象至指定存储桶
type objectPut struct {
    qcloud.BaseCos
    objectKey string // 对象名称
    filePath  string // 上传文件
}

func (op *objectPut) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        op.objectKey = "/" + objectKey
        op.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (op *objectPut) SetFilePath(filePath string) {
    f, err := os.Stat(filePath)
    if err != nil {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "文件不能是目录", nil))
    }
    op.filePath = filePath
}

func (op *objectPut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(op.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }
    if len(op.filePath) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传文件不能为空", nil))
    }

    // 新建一个缓冲，用于存放文件内容
    bodyBuffer := &bytes.Buffer{}
    // 创建一个multipart文件写入器,方便按照http规定格式写入内容
    bodyWriter := multipart.NewWriter(bodyBuffer)
    defer bodyWriter.Close()
    // 从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次
    fileWriter, err := bodyWriter.CreateFormFile("media", path.Base(op.filePath))
    if err != nil {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传文件失败", nil))
    }

    file, err := os.Open(op.filePath)
    defer file.Close()
    if err != nil {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传文件失败", nil))
    }

    _, err = io.Copy(fileWriter, file)
    if err != nil {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传文件失败", nil))
    }

    contentType := bodyWriter.FormDataContentType()
    content := bodyBuffer.Bytes()
    contentLength := len(string(content))

    encodeStr := base64.StdEncoding.EncodeToString(content)
    op.SetHeaderData("Content-MD5", mpf.HashMd5(encodeStr, ""))
    op.SetHeaderData("Content-Length", strconv.Itoa(contentLength))
    op.ReqUrl = "http://" + op.ReqHeader["Host"] + op.ReqUri
    client, req := op.GetRequest()
    req.Header.SetContentType(contentType)
    req.SetBody(content)

    return client, req
}

func NewObjectPut() *objectPut {
    op := &objectPut{qcloud.NewCos(), "", ""}
    op.ReqMethod = fasthttp.MethodPut
    return op
}
